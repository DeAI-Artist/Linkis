// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	big "math/big"
	time "time"

	mock "github.com/stretchr/testify/mock"

	types "github.com/DeAI-Artist/MintAI/common/types"
)

// Head is an autogenerated mock type for the Head type
type Head[BLOCK_HASH types.Hashable] struct {
	mock.Mock
}

// BlockDifficulty provides a mock function with given fields:
func (_m *Head[BLOCK_HASH]) BlockDifficulty() *big.Int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for BlockDifficulty")
	}

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// BlockHash provides a mock function with given fields:
func (_m *Head[BLOCK_HASH]) BlockHash() BLOCK_HASH {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for BlockHash")
	}

	var r0 BLOCK_HASH
	if rf, ok := ret.Get(0).(func() BLOCK_HASH); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(BLOCK_HASH)
	}

	return r0
}

// BlockNumber provides a mock function with given fields:
func (_m *Head[BLOCK_HASH]) BlockNumber() int64 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for BlockNumber")
	}

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// ChainLength provides a mock function with given fields:
func (_m *Head[BLOCK_HASH]) ChainLength() uint32 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ChainLength")
	}

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EarliestHeadInChain provides a mock function with given fields:
func (_m *Head[BLOCK_HASH]) EarliestHeadInChain() types.Head[BLOCK_HASH] {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for EarliestHeadInChain")
	}

	var r0 types.Head[BLOCK_HASH]
	if rf, ok := ret.Get(0).(func() types.Head[BLOCK_HASH]); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Head[BLOCK_HASH])
		}
	}

	return r0
}

// GetParent provides a mock function with given fields:
func (_m *Head[BLOCK_HASH]) GetParent() types.Head[BLOCK_HASH] {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetParent")
	}

	var r0 types.Head[BLOCK_HASH]
	if rf, ok := ret.Get(0).(func() types.Head[BLOCK_HASH]); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Head[BLOCK_HASH])
		}
	}

	return r0
}

// GetParentHash provides a mock function with given fields:
func (_m *Head[BLOCK_HASH]) GetParentHash() BLOCK_HASH {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetParentHash")
	}

	var r0 BLOCK_HASH
	if rf, ok := ret.Get(0).(func() BLOCK_HASH); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(BLOCK_HASH)
	}

	return r0
}

// GetTimestamp provides a mock function with given fields:
func (_m *Head[BLOCK_HASH]) GetTimestamp() time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetTimestamp")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// HashAtHeight provides a mock function with given fields: blockNum
func (_m *Head[BLOCK_HASH]) HashAtHeight(blockNum int64) BLOCK_HASH {
	ret := _m.Called(blockNum)

	if len(ret) == 0 {
		panic("no return value specified for HashAtHeight")
	}

	var r0 BLOCK_HASH
	if rf, ok := ret.Get(0).(func(int64) BLOCK_HASH); ok {
		r0 = rf(blockNum)
	} else {
		r0 = ret.Get(0).(BLOCK_HASH)
	}

	return r0
}

// NewHead creates a new instance of Head. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHead[BLOCK_HASH types.Hashable](t interface {
	mock.TestingT
	Cleanup(func())
}) *Head[BLOCK_HASH] {
	mock := &Head[BLOCK_HASH]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
