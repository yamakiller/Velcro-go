package retry

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/velcro-go/rpc/pkg/rpcinfo"
	"github.com/yamakiller/velcro-go/utils/circuitbreak"
	"github.com/yamakiller/velcro-go/utils/verrors"
	"github.com/yamakiller/velcro-go/vlog"
)

func newFailureRetryer(policy Policy, r *ShouldResultRetry, container *cbrContainer) (Retryer, error) {
	fr := &failureRetryer{specifiedResultRetry: r, container: container}
	if err := fr.UpdatePolicy(policy); err != nil {
		return nil, fmt.Errorf("newfailureRetryer failed, err=%w", err)
	}
	return fr, nil
}

type failureRetryer struct {
	enable               bool
	policy               *FailurePolicy
	backOff              BackOff
	container            *cbrContainer
	specifiedResultRetry *ShouldResultRetry
	sync.RWMutex
	errMsg string
}

// ShouldRetry implements the Retryer interface.
func (r *failureRetryer) ShouldRetry(ctx context.Context, err error, callTimes int, req interface{}, cbKey string) (string, bool) {
	r.RLock()
	defer r.RUnlock()
	if !r.enable {
		return "", false
	}
	if stop, msg := circuitBreakerStop(ctx, r.policy.StopPolicy, r.container, req, cbKey); stop {
		return msg, false
	}
	if stop, msg := ddlStop(ctx, r.policy.StopPolicy); stop {
		return msg, false
	}
	r.backOff.Wait(callTimes)
	return "", true
}

// AllowRetry implements the Retryer interface.
func (r *failureRetryer) AllowRetry(ctx context.Context) (string, bool) {
	r.RLock()
	defer r.RUnlock()
	if !r.enable || r.policy.StopPolicy.MaxRetryTimes == 0 {
		return "", false
	}
	if stop, msg := chainStop(ctx, r.policy.StopPolicy); stop {
		return msg, false
	}
	if stop, msg := ddlStop(ctx, r.policy.StopPolicy); stop {
		return msg, false
	}
	return "", true
}

// Do implement the Retryer interface.
func (r *failureRetryer) Do(ctx context.Context, rpcCall RPCCallFunc, firstRI rpcinfo.RPCInfo, req interface{}) (rpcinfo.RPCInfo, bool, error) {
	var err error
	r.RLock()
	var maxDuration time.Duration
	if r.policy.StopPolicy.MaxDurationMillisecond > 0 {
		maxDuration = time.Duration(r.policy.StopPolicy.MaxDurationMillisecond) * time.Millisecond
	}
	retryTimes := r.policy.StopPolicy.MaxRetryTimes
	r.RUnlock()

	var callTimes int32
	var callCosts strings.Builder
	var cRI rpcinfo.RPCInfo
	cbKey, _ := r.container.ctrl.GetKey(ctx, req)
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = panicToErr(ctx, panicInfo, firstRI)
		}
	}()
	startTime := time.Now()
	for i := 0; i <= retryTimes; i++ {
		var resp interface{}
		var callStart time.Time
		if i == 0 {
			callStart = startTime
		} else if i > 0 {
			if maxDuration > 0 && time.Since(startTime) > maxDuration {
				err = makeRetryErr(ctx, "exceed max duration", callTimes)
				break
			}
			if msg, ok := r.ShouldRetry(ctx, err, i, req, cbKey); !ok {
				if msg != "" {
					appendMsg := fmt.Sprintf("retried %d, %s", i-1, msg)
					appendErrMsg(err, appendMsg)
				}
				break
			}
			callStart = time.Now()
			callCosts.WriteByte(',')
			if respOp, ok := ctx.Value(ContextRespOp).(*int32); ok {
				atomic.StoreInt32(respOp, OpNo)
			}
		}
		callTimes++
		if r.container.enablePercentageLimit {
			// record stat before call since requests may be slow, making the limiter more accurate
			recordRetryStat(cbKey, r.container.panel, callTimes)
		}
		cRI, resp, err = rpcCall(ctx, r)
		callCosts.WriteString(strconv.FormatInt(time.Since(callStart).Microseconds(), 10))

		if !r.container.enablePercentageLimit && r.container.stat {
			circuitbreak.RecordStat(ctx, req, nil, err, cbKey, r.container.ctrl, r.container.panel)
		}
		if err == nil {
			if r.policy.IsRespRetryNonNil() && r.policy.ShouldResultRetry.RespRetry(resp, cRI) {
				// user specified resp to do retry
				continue
			}
			break
		} else {
			if i == retryTimes {
				// stop retry then wrap error
				err = verrors.ErrRetry.WithCause(err)
			} else if !r.isRetryErr(err, cRI) {
				// not timeout or user specified error won't do retry
				break
			}
		}
	}
	recordRetryInfo(cRI, callTimes, callCosts.String())
	if err == nil && callTimes == 1 {
		return cRI, true, nil
	}
	return cRI, false, err
}

// UpdatePolicy implements the Retryer interface.
func (r *failureRetryer) UpdatePolicy(rp Policy) (err error) {
	if !rp.Enable {
		r.Lock()
		r.enable = rp.Enable
		r.Unlock()
		return nil
	}
	var errMsg string
	if rp.FailurePolicy == nil || rp.Type != FailureType {
		errMsg = "FailurePolicy is nil or retry type not match, cannot do update in failureRetryer"
		err = errors.New(errMsg)
	}
	rt := rp.FailurePolicy.StopPolicy.MaxRetryTimes
	if errMsg == "" && (rt < 0 || rt > maxFailureRetryTimes) {
		errMsg = fmt.Sprintf("invalid failure MaxRetryTimes[%d]", rt)
		err = errors.New(errMsg)
	}
	if errMsg == "" {
		if e := checkCircuitBreakerErrorRate(&rp.FailurePolicy.StopPolicy.CircuitBreakPolicy); e != nil {
			rp.FailurePolicy.StopPolicy.CircuitBreakPolicy.ErrorRate = defaultCircuitBreakerErrRate
			errMsg = fmt.Sprintf("failureRetryer %s, use default %0.2f", e.Error(), defaultCircuitBreakerErrRate)
			vlog.Warnf(errMsg)
		}
	}
	r.Lock()
	defer r.Unlock()
	r.enable = rp.Enable
	if err != nil {
		r.errMsg = errMsg
		return err
	}
	r.policy = rp.FailurePolicy
	r.setSpecifiedResultRetryIfNeeded(r.specifiedResultRetry)
	if bo, e := initBackOff(rp.FailurePolicy.BackOffPolicy); e != nil {
		r.errMsg = fmt.Sprintf("failureRetryer update BackOffPolicy failed, err=%s", e.Error())
		vlog.Warnf(r.errMsg)
	} else {
		r.backOff = bo
	}
	return nil
}

// AppendErrMsgIfNeeded implements the Retryer interface.
func (r *failureRetryer) AppendErrMsgIfNeeded(err error, ri rpcinfo.RPCInfo, msg string) {
	if r.isRetryErr(err, ri) {
		// Add additional reason when retry is not applied.
		appendErrMsg(err, msg)
	}
}

// Prepare implements the Retryer interface.
func (r *failureRetryer) Prepare(ctx context.Context, prevRI, retryRI rpcinfo.RPCInfo) {
	handleRetryInstance(r.policy.RetrySameNode, prevRI, retryRI)
}

func (r *failureRetryer) isRetryErr(err error, ri rpcinfo.RPCInfo) bool {
	if err == nil {
		return false
	}
	// Logic Notice:
	// some kinds of error cannot be retried, eg: ServiceCircuitBreak.
	// But CircuitBreak has been checked in ShouldRetry, it doesn't need to filter ServiceCircuitBreak.
	// If there are some other specified errors that cannot be retried, it should be filtered here.

	if r.policy.IsRetryForTimeout() && verrors.IsTimeoutError(err) {
		return true
	}
	if r.policy.IsErrorRetryNonNil() && r.policy.ShouldResultRetry.ErrorRetry(err, ri) {
		return true
	}
	return false
}

func initBackOff(policy *BackOffPolicy) (bo BackOff, err error) {
	bo = NoneBackOff
	if policy == nil {
		return
	}
	switch policy.BackOffType {
	case NoneBackOffType:
	case FixedBackOffType:
		if policy.CfgItems == nil {
			return bo, errors.New("invalid FixedBackOff, CfgItems is nil")
		}
		fixMS := policy.CfgItems[FixMSBackOffCfgKey]
		fixMSInt, _ := strconv.Atoi(fmt.Sprintf("%1.0f", fixMS))
		if err = checkFixedBackOff(fixMSInt); err != nil {
			return
		}
		bo = newFixedBackOff(fixMSInt)
	case RandomBackOffType:
		if policy.CfgItems == nil {
			return bo, errors.New("invalid FixedBackOff, CfgItems is nil")
		}
		minMS := policy.CfgItems[MinMSBackOffCfgKey]
		maxMS := policy.CfgItems[MaxMSBackOffCfgKey]
		minMSInt, _ := strconv.Atoi(fmt.Sprintf("%1.0f", minMS))
		maxMSInt, _ := strconv.Atoi(fmt.Sprintf("%1.0f", maxMS))
		if err = checkRandomBackOff(minMSInt, maxMSInt); err != nil {
			return
		}
		bo = newRandomBackOff(minMSInt, maxMSInt)
	default:
		return bo, fmt.Errorf("invalid backoffType=%v", policy.BackOffType)
	}
	return
}

// Type implements the Retryer interface.
func (r *failureRetryer) Type() Type {
	return FailureType
}

// Dump implements the Retryer interface.
func (r *failureRetryer) Dump() map[string]interface{} {
	r.RLock()
	defer r.RUnlock()
	dm := make(map[string]interface{})
	dm["enable"] = r.enable
	dm["failure_retry"] = r.policy
	if r.policy != nil {
		dm["specified_result_retry"] = map[string]bool{
			"error_retry": r.policy.IsErrorRetryNonNil(),
			"resp_retry":  r.policy.IsRespRetryNonNil(),
		}
	}
	if r.errMsg != "" {
		dm["errMsg"] = r.errMsg
	}
	return dm
}

func (r *failureRetryer) setSpecifiedResultRetryIfNeeded(rr *ShouldResultRetry) {
	if rr != nil {
		// save the object specified by client.WithSpecifiedResultRetry(..)
		r.specifiedResultRetry = rr
	}
	if r.policy != nil && r.specifiedResultRetry != nil {
		// The priority of client.WithSpecifiedResultRetry(..) is higher, so always update it
		// NOTE: client.WithSpecifiedResultRetry(..) will always reject a nil object
		r.policy.ShouldResultRetry = r.specifiedResultRetry
	}
}
