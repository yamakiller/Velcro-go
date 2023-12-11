package containers

import (
	"sync"

	"github.com/yamakiller/velcro-go/utils"
)

func NewQueue(cap int) *Queue {
	return &Queue{
		_cap:               cap,
		_overloadThreshold: cap * 2,
		_buffer:            make([]interface{}, cap),
	}
}

// Queue 连续可膨胀队列
type Queue struct {
	_cap  int
	_head int
	_tail int

	_overload          int
	_overloadThreshold int

	_buffer []interface{}
	_sync   sync.Mutex
}

// Push Insert an object
// @Param (interface{}) item
func (qe *Queue) Push(item interface{}) {
	qe._sync.Lock()
	defer qe._sync.Unlock()
	qe.unpush(item)
}

// Pop doc
// @Method Pop @Summary Take an object, If empty return nil
// @Return (interface{}) return object
// @Return (bool)
func (qe *Queue) Pop() (interface{}, bool) {
	qe._sync.Lock()
	defer qe._sync.Unlock()
	return qe.unpop()
}

// Overload Detecting queues exceeding the limit [mainly used for warning records]
// @Return (int)
func (qe *Queue) Overload() int {
	if qe._overload != 0 {
		overload := qe._overload
		qe._overload = 0
		return overload
	}
	return 0
}

// Length Length of the Queue
// @Return (int) length
func (qe *Queue) Length() int {
	var (
		head int
		tail int
		cap  int
	)
	qe._sync.Lock()
	head = qe._head
	tail = qe._tail
	cap = qe._cap
	qe._sync.Unlock()

	if head <= tail {
		return tail - head
	}
	return tail + cap - head
}

func (qe *Queue) unpush(item interface{}) {
	utils.AssertEmpty(item, "error push is nil")
	qe._buffer[qe._tail] = item
	qe._tail++
	if qe._tail >= qe._cap {
		qe._tail = 0
	}

	if qe._head == qe._tail {
		qe.expand()
	}
}

func (qe *Queue) unpop() (interface{}, bool) {
	var resultSucces bool
	var result interface{}
	if qe._head != qe._tail {
		resultSucces = true
		result = qe._buffer[qe._head]
		qe._buffer[qe._head] = nil
		qe._head++
		if qe._head >= qe._cap {
			qe._head = 0
		}

		length := qe._tail - qe._tail
		if length < 0 {
			length += qe._cap
		}
		for length > qe._overloadThreshold {
			qe._overload = length
			qe._overloadThreshold *= 2
		}
	}

	return result, resultSucces
}

func (qe *Queue) expand() {
	newBuff := make([]interface{}, qe._cap*2)
	for i := 0; i < qe._cap; i++ {
		newBuff[i] = qe._buffer[(qe._head+i)%qe._cap]
	}

	qe._head = 0
	qe._tail = qe._cap
	qe._cap *= 2

	qe._buffer = newBuff
}
