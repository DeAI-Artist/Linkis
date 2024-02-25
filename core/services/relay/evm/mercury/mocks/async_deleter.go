// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	pb "github.com/DeAI-Artist/MintAI/core/services/relay/evm/mercury/wsrpc/pb"
	mock "github.com/stretchr/testify/mock"
)

// AsyncDeleter is an autogenerated mock type for the asyncDeleter type
type AsyncDeleter struct {
	mock.Mock
}

// AsyncDelete provides a mock function with given fields: req
func (_m *AsyncDeleter) AsyncDelete(req *pb.TransmitRequest) {
	_m.Called(req)
}

// NewAsyncDeleter creates a new instance of AsyncDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAsyncDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *AsyncDeleter {
	mock := &AsyncDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
