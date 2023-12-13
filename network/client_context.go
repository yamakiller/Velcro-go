package network

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/yamakiller/velcro-go/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// 通用客户端上下文
type clientContext struct {
	_client      Client
	_system      *NetworkSystem
	_self        *ClientID
	_message     []byte
	_messageFrom *net.Addr
	_state       int32
}

func (ctx *clientContext) Client() Client {
	return ctx._client
}

func (ctx *clientContext) Self() *ClientID {
	return ctx._self
}

func (ctx *clientContext) NetworkSystem() *NetworkSystem {
	return ctx._system
}

// Message 返回当前接收到的消息
func (ctx *clientContext) Message() []byte {
	return ctx._message
}

func (ctx *clientContext) MessageFrom() *net.Addr {
	return ctx._messageFrom
}

func (ctx *clientContext) PostMessage(cid *ClientID, message []byte) {

	systemMetrics, ok := ctx._system._extensions.Get(extensionId).(*Metrics)
	if ok && systemMetrics._enabled {
		t := time.Now()
		cid.postMessage(ctx._system, message)

		delta := time.Since(t)
		_ctx := context.Background()

		if instruments := systemMetrics._metrics.Get(metrics.InternalClientMetrics); instruments != nil {
			histogram := instruments.ClientBytesSendHistogram

			labels := append(
				systemMetrics.CommonLabels(ctx),
				attribute.String("message bytes", fmt.Sprintf("%d", len(message))),
			)
			histogram.Record(_ctx, delta.Seconds(), metric.WithAttributes(labels...))
		}
	} else {
		cid.postMessage(ctx._system, message)
	}
}

func (ctx *clientContext) PostToMessage(cid *ClientID, message []byte, target *net.Addr) {
	cid.postToMessage(ctx._system, message, *target)
}

// Close 关闭当前 Client
func (ctx *clientContext) Close() {
	if atomic.LoadInt32(&ctx._state) >= stateClosing {
		// 已经关闭
		return
	}

	if ctx._system.Config.MetricsProvider != nil {
		metricsSystem, ok := ctx._system._extensions.Get(extensionId).(*Metrics)
		if ok && metricsSystem._enabled {
			_ctx := context.Background()
			if instruments := metricsSystem._metrics.Get(metrics.InternalClientMetrics); instruments != nil {
				instruments.ClientCloseCount.Add(_ctx, 1, metric.WithAttributes(metricsSystem.CommonLabels(ctx)...))
			}
		}
	}

	atomic.StoreInt32(&ctx._state, stateClosing)
	ctx._self.Close(ctx._system)
}

func (ctx *clientContext) incarnateClient() {
	atomic.StoreInt32(&ctx._state, stateAccept)

	ctx._client = ctx._system._producer(ctx._system)

	if ctx._system.Config.MetricsProvider != nil {
		metricsSystem, ok := ctx._system._extensions.Get(extensionId).(*Metrics)
		if ok && metricsSystem._enabled {
			_ctx := context.Background()
			if instruments := metricsSystem._metrics.Get(metrics.InternalClientMetrics); instruments != nil {
				instruments.ClientSpawnCount.Add(_ctx, 1, metric.WithAttributes(metricsSystem.CommonLabels(ctx)...))
			}
		}
	}
}

// 日志接口
// Info ...
func (ctx *clientContext) Info(sfmt string, args ...interface{}) {
	ctx._system._logger.Info(fmt.Sprintf("[%s]", ctx._self.ToString()), sfmt, args...)
}

// Debug ...
func (ctx *clientContext) Debug(sfmt string, args ...interface{}) {
	ctx._system._logger.Debug(fmt.Sprintf("[%s]", ctx._self.ToString()), sfmt, args...)
}

// Error ...
func (ctx *clientContext) Error(sfmt string, args ...interface{}) {
	ctx._system._logger.Error(fmt.Sprintf("[%s]", ctx._self.ToString()), sfmt, args...)
}

// Warning ...
func (ctx *clientContext) Warning(sfmt string, args ...interface{}) {
	ctx._system._logger.Warning(fmt.Sprintf("[%s]", ctx._self.ToString()), sfmt, args...)
}

// Fatal ...
func (ctx *clientContext) Fatal(sfmt string, args ...interface{}) {
	ctx._system._logger.Fatal(fmt.Sprintf("[%s]", ctx._self.ToString()), sfmt, args...)
}

// Panic ...
func (ctx *clientContext) Panic(sfmt string, args ...interface{}) {
	ctx._system._logger.Panic(fmt.Sprintf("[%s]", ctx._self.ToString()), sfmt, args...)
}

func (ctx *clientContext) invokerAccept() {
	ctx._client.Accept(ctx)
}

// MessageInvoker 接口实现
func (ctx *clientContext) invokerRecvice(b []byte, addr *net.Addr) {
	if atomic.LoadInt32(&ctx._state) == stateClosed {
		// 已经关闭
		return
	}

	systemMetrics, ok := ctx._system._extensions.Get(extensionId).(*Metrics)
	if ok && systemMetrics._enabled {
		t := time.Now()
		ctx._message = b
		ctx._messageFrom = addr
		ctx._client.Recvice(ctx)
		ctx._messageFrom = nil
		ctx._message = nil

		delta := time.Since(t)
		_ctx := context.Background()

		if instruments := systemMetrics._metrics.Get(metrics.InternalClientMetrics); instruments != nil {
			histogram := instruments.ClientBytesRecviceHistogram

			labels := append(
				systemMetrics.CommonLabels(ctx),
				attribute.String("message bytes", fmt.Sprintf("%d", len(b))),
			)
			histogram.Record(_ctx, delta.Seconds(), metric.WithAttributes(labels...))
		}
	} else {
		ctx._message = b
		ctx._messageFrom = addr
		ctx._client.Recvice(ctx)
		ctx._messageFrom = nil
		ctx._message = nil
	}

}

func (ctx *clientContext) invokerClosed() {
	if atomic.LoadInt32(&ctx._state) == stateClosed {
		// 已关闭
		return
	}
	ctx._system._handlers.Remove(ctx._self)
	atomic.StoreInt32(&ctx._state, stateClosed)
	// 调用客户端逻辑
	ctx._client.Closed(ctx)
	// TODO: 释放Client
}
