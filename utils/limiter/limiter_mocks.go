package limiter

import (
	"context"
	"reflect"
	"time"

	"github.com/golang/mock/gomock"
)

// MockConcurrencyLimiter is a mock of ConcurrencyLimiter interface.
type MockConcurrencyLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockConcurrencyLimiterMockRecorder
}

// MockConcurrencyLimiterMockRecorder 是 MockConcurrencyLimiter 的模拟记录器.
type MockConcurrencyLimiterMockRecorder struct {
	mock *MockConcurrencyLimiter
}

// NewMockConcurrencyLimiter 创建一个新的模拟实例.
func NewMockConcurrencyLimiter(ctrl *gomock.Controller) *MockConcurrencyLimiter {
	mock := &MockConcurrencyLimiter{ctrl: ctrl}
	mock.recorder = &MockConcurrencyLimiterMockRecorder{mock}
	return mock
}

// EXPECT 返回一个对象，允许调用者指示预期用途.
func (m *MockConcurrencyLimiter) EXPECT() *MockConcurrencyLimiterMockRecorder {
	return m.recorder
}

// Acquire mocks 基础方法.
func (m *MockConcurrencyLimiter) Acquire(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Acquire", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Acquire 表示预期的 Acquire 调用.
func (mr *MockConcurrencyLimiterMockRecorder) Acquire(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Acquire", reflect.TypeOf((*MockConcurrencyLimiter)(nil).Acquire), ctx)
}

// Release mocks 基础方法.
func (m *MockConcurrencyLimiter) Release(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Release", ctx)
}

// Release 表示预期的 Release 调用.
func (mr *MockConcurrencyLimiterMockRecorder) Release(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*MockConcurrencyLimiter)(nil).Release), ctx)
}

// Status mocks 基础方法.
func (m *MockConcurrencyLimiter) Status(ctx context.Context) (int, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// Status 指示 Status 的预期调用.
func (mr *MockConcurrencyLimiterMockRecorder) Status(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockConcurrencyLimiter)(nil).Status), ctx)
}

// MockRateLimiter 是 RateLimiter 接口的模拟.
type MockRateLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockRateLimiterMockRecorder
}

// MockRateLimiterMockRecorder 是 MockRateLimiter 的模拟记录器.
type MockRateLimiterMockRecorder struct {
	mock *MockRateLimiter
}

// NewMockRateLimiter 创建一个新的模拟实例.
func NewMockRateLimiter(ctrl *gomock.Controller) *MockRateLimiter {
	mock := &MockRateLimiter{ctrl: ctrl}
	mock.recorder = &MockRateLimiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateLimiter) EXPECT() *MockRateLimiterMockRecorder {
	return m.recorder
}

// Acquire mocks base method.
func (m *MockRateLimiter) Acquire(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Acquire", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Acquire indicates an expected call of Acquire.
func (mr *MockRateLimiterMockRecorder) Acquire(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Acquire", reflect.TypeOf((*MockRateLimiter)(nil).Acquire), ctx)
}

// Status mocks base method.
func (m *MockRateLimiter) Status(ctx context.Context) (int, int, time.Duration) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(time.Duration)
	return ret0, ret1, ret2
}

// Status indicates an expected call of Status.
func (mr *MockRateLimiterMockRecorder) Status(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockRateLimiter)(nil).Status), ctx)
}

// MockUpdatable is a mock of Updatable interface.
type MockUpdatable struct {
	ctrl     *gomock.Controller
	recorder *MockUpdatableMockRecorder
}

// MockUpdatableMockRecorder is the mock recorder for MockUpdatable.
type MockUpdatableMockRecorder struct {
	mock *MockUpdatable
}

// NewMockUpdatable creates a new mock instance.
func NewMockUpdatable(ctrl *gomock.Controller) *MockUpdatable {
	mock := &MockUpdatable{ctrl: ctrl}
	mock.recorder = &MockUpdatableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpdatable) EXPECT() *MockUpdatableMockRecorder {
	return m.recorder
}

// UpdateLimit mocks base method.
func (m *MockUpdatable) UpdateLimit(limit int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateLimit", limit)
}

// UpdateLimit indicates an expected call of UpdateLimit.
func (mr *MockUpdatableMockRecorder) UpdateLimit(limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLimit", reflect.TypeOf((*MockUpdatable)(nil).UpdateLimit), limit)
}

// MockLimitReporter is a mock of LimitReporter interface.
type MockLimitReporter struct {
	ctrl     *gomock.Controller
	recorder *MockLimitReporterMockRecorder
}

// MockLimitReporterMockRecorder is the mock recorder for MockLimitReporter.
type MockLimitReporterMockRecorder struct {
	mock *MockLimitReporter
}

// NewMockLimitReporter creates a new mock instance.
func NewMockLimitReporter(ctrl *gomock.Controller) *MockLimitReporter {
	mock := &MockLimitReporter{ctrl: ctrl}
	mock.recorder = &MockLimitReporterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLimitReporter) EXPECT() *MockLimitReporterMockRecorder {
	return m.recorder
}

// ConnOverloadReport mocks base method.
func (m *MockLimitReporter) ConnOverloadReport() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ConnOverloadReport")
}

// ConnOverloadReport indicates an expected call of ConnOverloadReport.
func (mr *MockLimitReporterMockRecorder) ConnOverloadReport() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnOverloadReport", reflect.TypeOf((*MockLimitReporter)(nil).ConnOverloadReport))
}

// QPSOverloadReport mocks base method.
func (m *MockLimitReporter) QPSOverloadReport() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "QPSOverloadReport")
}

// QPSOverloadReport indicates an expected call of QPSOverloadReport.
func (mr *MockLimitReporterMockRecorder) QPSOverloadReport() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QPSOverloadReport", reflect.TypeOf((*MockLimitReporter)(nil).QPSOverloadReport))
}
