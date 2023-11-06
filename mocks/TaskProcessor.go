// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	asynq "github.com/hibiken/asynq"

	mock "github.com/stretchr/testify/mock"
)

// TaskProcessor is an autogenerated mock type for the TaskProcessor type
type TaskProcessor struct {
	mock.Mock
}

type TaskProcessor_Expecter struct {
	mock *mock.Mock
}

func (_m *TaskProcessor) EXPECT() *TaskProcessor_Expecter {
	return &TaskProcessor_Expecter{mock: &_m.Mock}
}

// ProcessTaskSendVerifyEmail provides a mock function with given fields: ctx, task
func (_m *TaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	ret := _m.Called(ctx, task)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *asynq.Task) error); ok {
		r0 = rf(ctx, task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TaskProcessor_ProcessTaskSendVerifyEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProcessTaskSendVerifyEmail'
type TaskProcessor_ProcessTaskSendVerifyEmail_Call struct {
	*mock.Call
}

// ProcessTaskSendVerifyEmail is a helper method to define mock.On call
//  - ctx context.Context
//  - task *asynq.Task
func (_e *TaskProcessor_Expecter) ProcessTaskSendVerifyEmail(ctx interface{}, task interface{}) *TaskProcessor_ProcessTaskSendVerifyEmail_Call {
	return &TaskProcessor_ProcessTaskSendVerifyEmail_Call{Call: _e.mock.On("ProcessTaskSendVerifyEmail", ctx, task)}
}

func (_c *TaskProcessor_ProcessTaskSendVerifyEmail_Call) Run(run func(ctx context.Context, task *asynq.Task)) *TaskProcessor_ProcessTaskSendVerifyEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*asynq.Task))
	})
	return _c
}

func (_c *TaskProcessor_ProcessTaskSendVerifyEmail_Call) Return(_a0 error) *TaskProcessor_ProcessTaskSendVerifyEmail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TaskProcessor_ProcessTaskSendVerifyEmail_Call) RunAndReturn(run func(context.Context, *asynq.Task) error) *TaskProcessor_ProcessTaskSendVerifyEmail_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields:
func (_m *TaskProcessor) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TaskProcessor_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type TaskProcessor_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
func (_e *TaskProcessor_Expecter) Start() *TaskProcessor_Start_Call {
	return &TaskProcessor_Start_Call{Call: _e.mock.On("Start")}
}

func (_c *TaskProcessor_Start_Call) Run(run func()) *TaskProcessor_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TaskProcessor_Start_Call) Return(_a0 error) *TaskProcessor_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TaskProcessor_Start_Call) RunAndReturn(run func() error) *TaskProcessor_Start_Call {
	_c.Call.Return(run)
	return _c
}

// NewTaskProcessor creates a new instance of TaskProcessor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskProcessor(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskProcessor {
	mock := &TaskProcessor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
