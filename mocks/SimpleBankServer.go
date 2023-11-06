// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	pb "github.com/tgfukuda/be-master/pb"
)

// SimpleBankServer is an autogenerated mock type for the SimpleBankServer type
type SimpleBankServer struct {
	mock.Mock
}

type SimpleBankServer_Expecter struct {
	mock *mock.Mock
}

func (_m *SimpleBankServer) EXPECT() *SimpleBankServer_Expecter {
	return &SimpleBankServer_Expecter{mock: &_m.Mock}
}

// CreateUser provides a mock function with given fields: _a0, _a1
func (_m *SimpleBankServer) CreateUser(_a0 context.Context, _a1 *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *pb.CreateUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.CreateUserRequest) (*pb.CreateUserResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pb.CreateUserRequest) *pb.CreateUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.CreateUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pb.CreateUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SimpleBankServer_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type SimpleBankServer_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//  - _a0 context.Context
//  - _a1 *pb.CreateUserRequest
func (_e *SimpleBankServer_Expecter) CreateUser(_a0 interface{}, _a1 interface{}) *SimpleBankServer_CreateUser_Call {
	return &SimpleBankServer_CreateUser_Call{Call: _e.mock.On("CreateUser", _a0, _a1)}
}

func (_c *SimpleBankServer_CreateUser_Call) Run(run func(_a0 context.Context, _a1 *pb.CreateUserRequest)) *SimpleBankServer_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*pb.CreateUserRequest))
	})
	return _c
}

func (_c *SimpleBankServer_CreateUser_Call) Return(_a0 *pb.CreateUserResponse, _a1 error) *SimpleBankServer_CreateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SimpleBankServer_CreateUser_Call) RunAndReturn(run func(context.Context, *pb.CreateUserRequest) (*pb.CreateUserResponse, error)) *SimpleBankServer_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// LoginUser provides a mock function with given fields: _a0, _a1
func (_m *SimpleBankServer) LoginUser(_a0 context.Context, _a1 *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *pb.LoginUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.LoginUserRequest) (*pb.LoginUserResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pb.LoginUserRequest) *pb.LoginUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.LoginUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pb.LoginUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SimpleBankServer_LoginUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoginUser'
type SimpleBankServer_LoginUser_Call struct {
	*mock.Call
}

// LoginUser is a helper method to define mock.On call
//  - _a0 context.Context
//  - _a1 *pb.LoginUserRequest
func (_e *SimpleBankServer_Expecter) LoginUser(_a0 interface{}, _a1 interface{}) *SimpleBankServer_LoginUser_Call {
	return &SimpleBankServer_LoginUser_Call{Call: _e.mock.On("LoginUser", _a0, _a1)}
}

func (_c *SimpleBankServer_LoginUser_Call) Run(run func(_a0 context.Context, _a1 *pb.LoginUserRequest)) *SimpleBankServer_LoginUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*pb.LoginUserRequest))
	})
	return _c
}

func (_c *SimpleBankServer_LoginUser_Call) Return(_a0 *pb.LoginUserResponse, _a1 error) *SimpleBankServer_LoginUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SimpleBankServer_LoginUser_Call) RunAndReturn(run func(context.Context, *pb.LoginUserRequest) (*pb.LoginUserResponse, error)) *SimpleBankServer_LoginUser_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUser provides a mock function with given fields: _a0, _a1
func (_m *SimpleBankServer) UpdateUser(_a0 context.Context, _a1 *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *pb.UpdateUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pb.UpdateUserRequest) *pb.UpdateUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.UpdateUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pb.UpdateUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SimpleBankServer_UpdateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUser'
type SimpleBankServer_UpdateUser_Call struct {
	*mock.Call
}

// UpdateUser is a helper method to define mock.On call
//  - _a0 context.Context
//  - _a1 *pb.UpdateUserRequest
func (_e *SimpleBankServer_Expecter) UpdateUser(_a0 interface{}, _a1 interface{}) *SimpleBankServer_UpdateUser_Call {
	return &SimpleBankServer_UpdateUser_Call{Call: _e.mock.On("UpdateUser", _a0, _a1)}
}

func (_c *SimpleBankServer_UpdateUser_Call) Run(run func(_a0 context.Context, _a1 *pb.UpdateUserRequest)) *SimpleBankServer_UpdateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*pb.UpdateUserRequest))
	})
	return _c
}

func (_c *SimpleBankServer_UpdateUser_Call) Return(_a0 *pb.UpdateUserResponse, _a1 error) *SimpleBankServer_UpdateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SimpleBankServer_UpdateUser_Call) RunAndReturn(run func(context.Context, *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)) *SimpleBankServer_UpdateUser_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyEmail provides a mock function with given fields: _a0, _a1
func (_m *SimpleBankServer) VerifyEmail(_a0 context.Context, _a1 *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *pb.VerifyEmailResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pb.VerifyEmailRequest) *pb.VerifyEmailResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.VerifyEmailResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pb.VerifyEmailRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SimpleBankServer_VerifyEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyEmail'
type SimpleBankServer_VerifyEmail_Call struct {
	*mock.Call
}

// VerifyEmail is a helper method to define mock.On call
//  - _a0 context.Context
//  - _a1 *pb.VerifyEmailRequest
func (_e *SimpleBankServer_Expecter) VerifyEmail(_a0 interface{}, _a1 interface{}) *SimpleBankServer_VerifyEmail_Call {
	return &SimpleBankServer_VerifyEmail_Call{Call: _e.mock.On("VerifyEmail", _a0, _a1)}
}

func (_c *SimpleBankServer_VerifyEmail_Call) Run(run func(_a0 context.Context, _a1 *pb.VerifyEmailRequest)) *SimpleBankServer_VerifyEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*pb.VerifyEmailRequest))
	})
	return _c
}

func (_c *SimpleBankServer_VerifyEmail_Call) Return(_a0 *pb.VerifyEmailResponse, _a1 error) *SimpleBankServer_VerifyEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SimpleBankServer_VerifyEmail_Call) RunAndReturn(run func(context.Context, *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error)) *SimpleBankServer_VerifyEmail_Call {
	_c.Call.Return(run)
	return _c
}

// mustEmbedUnimplementedSimpleBankServer provides a mock function with given fields:
func (_m *SimpleBankServer) mustEmbedUnimplementedSimpleBankServer() {
	_m.Called()
}

// SimpleBankServer_mustEmbedUnimplementedSimpleBankServer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'mustEmbedUnimplementedSimpleBankServer'
type SimpleBankServer_mustEmbedUnimplementedSimpleBankServer_Call struct {
	*mock.Call
}

// mustEmbedUnimplementedSimpleBankServer is a helper method to define mock.On call
func (_e *SimpleBankServer_Expecter) mustEmbedUnimplementedSimpleBankServer() *SimpleBankServer_mustEmbedUnimplementedSimpleBankServer_Call {
	return &SimpleBankServer_mustEmbedUnimplementedSimpleBankServer_Call{Call: _e.mock.On("mustEmbedUnimplementedSimpleBankServer")}
}

func (_c *SimpleBankServer_mustEmbedUnimplementedSimpleBankServer_Call) Run(run func()) *SimpleBankServer_mustEmbedUnimplementedSimpleBankServer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SimpleBankServer_mustEmbedUnimplementedSimpleBankServer_Call) Return() *SimpleBankServer_mustEmbedUnimplementedSimpleBankServer_Call {
	_c.Call.Return()
	return _c
}

func (_c *SimpleBankServer_mustEmbedUnimplementedSimpleBankServer_Call) RunAndReturn(run func()) *SimpleBankServer_mustEmbedUnimplementedSimpleBankServer_Call {
	_c.Call.Return(run)
	return _c
}

// NewSimpleBankServer creates a new instance of SimpleBankServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSimpleBankServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *SimpleBankServer {
	mock := &SimpleBankServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}