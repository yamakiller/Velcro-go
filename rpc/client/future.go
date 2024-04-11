package client

import (
	"sync"
	"time"

	"github.com/yamakiller/velcro-go/rpc/messages"
)

// 请求器
type Future struct {
	sequenceID int32
	cond       *sync.Cond
	done       bool
	request    *messages.RpcRequestMessage
	result     []byte
	err        error
	t          *time.Timer
}

// Error 错误信息
func (ref *Future) Error() error {
	return ref.err
}

// Result 结果
func (ref *Future) Result() []byte {
	return ref.result
}

func (ref *Future) Wait() {
	ref.cond.L.Lock()
	for !ref.done {
		ref.cond.Wait()
	}
	ref.cond.L.Unlock()
}
