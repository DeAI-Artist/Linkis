// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	big "github.com/DeAI-Artist/MintAI/core/chains/evm/utils/big"
	mock "github.com/stretchr/testify/mock"

	pg "github.com/DeAI-Artist/MintAI/core/services/pg"

	s4 "github.com/DeAI-Artist/MintAI/core/services/s4"

	time "time"
)

// ORM is an autogenerated mock type for the ORM type
type ORM struct {
	mock.Mock
}

// DeleteExpired provides a mock function with given fields: limit, utcNow, qopts
func (_m *ORM) DeleteExpired(limit uint, utcNow time.Time, qopts ...pg.QOpt) (int64, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, limit, utcNow)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteExpired")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(uint, time.Time, ...pg.QOpt) (int64, error)); ok {
		return rf(limit, utcNow, qopts...)
	}
	if rf, ok := ret.Get(0).(func(uint, time.Time, ...pg.QOpt) int64); ok {
		r0 = rf(limit, utcNow, qopts...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(uint, time.Time, ...pg.QOpt) error); ok {
		r1 = rf(limit, utcNow, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: address, slotId, qopts
func (_m *ORM) Get(address *big.Big, slotId uint, qopts ...pg.QOpt) (*s4.Row, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, address, slotId)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *s4.Row
	var r1 error
	if rf, ok := ret.Get(0).(func(*big.Big, uint, ...pg.QOpt) (*s4.Row, error)); ok {
		return rf(address, slotId, qopts...)
	}
	if rf, ok := ret.Get(0).(func(*big.Big, uint, ...pg.QOpt) *s4.Row); ok {
		r0 = rf(address, slotId, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s4.Row)
		}
	}

	if rf, ok := ret.Get(1).(func(*big.Big, uint, ...pg.QOpt) error); ok {
		r1 = rf(address, slotId, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSnapshot provides a mock function with given fields: addressRange, qopts
func (_m *ORM) GetSnapshot(addressRange *s4.AddressRange, qopts ...pg.QOpt) ([]*s4.SnapshotRow, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, addressRange)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetSnapshot")
	}

	var r0 []*s4.SnapshotRow
	var r1 error
	if rf, ok := ret.Get(0).(func(*s4.AddressRange, ...pg.QOpt) ([]*s4.SnapshotRow, error)); ok {
		return rf(addressRange, qopts...)
	}
	if rf, ok := ret.Get(0).(func(*s4.AddressRange, ...pg.QOpt) []*s4.SnapshotRow); ok {
		r0 = rf(addressRange, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*s4.SnapshotRow)
		}
	}

	if rf, ok := ret.Get(1).(func(*s4.AddressRange, ...pg.QOpt) error); ok {
		r1 = rf(addressRange, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUnconfirmedRows provides a mock function with given fields: limit, qopts
func (_m *ORM) GetUnconfirmedRows(limit uint, qopts ...pg.QOpt) ([]*s4.Row, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, limit)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetUnconfirmedRows")
	}

	var r0 []*s4.Row
	var r1 error
	if rf, ok := ret.Get(0).(func(uint, ...pg.QOpt) ([]*s4.Row, error)); ok {
		return rf(limit, qopts...)
	}
	if rf, ok := ret.Get(0).(func(uint, ...pg.QOpt) []*s4.Row); ok {
		r0 = rf(limit, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*s4.Row)
		}
	}

	if rf, ok := ret.Get(1).(func(uint, ...pg.QOpt) error); ok {
		r1 = rf(limit, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: row, qopts
func (_m *ORM) Update(row *s4.Row, qopts ...pg.QOpt) error {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, row)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*s4.Row, ...pg.QOpt) error); ok {
		r0 = rf(row, qopts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewORM creates a new instance of ORM. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewORM(t interface {
	mock.TestingT
	Cleanup(func())
}) *ORM {
	mock := &ORM{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
