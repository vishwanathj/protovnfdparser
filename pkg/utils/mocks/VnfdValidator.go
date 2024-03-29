// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// VnfdValidator is an autogenerated mock type for the VnfdValidator type
type VnfdValidator struct {
	mock.Mock
}

// ValidatePaginatedVnfdsInstancesBody provides a mock function with given fields: _a0
func (_m *VnfdValidator) ValidatePaginatedVnfdsInstancesBody(_a0 []byte) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateVnfdInstanceBody provides a mock function with given fields: _a0
func (_m *VnfdValidator) ValidateVnfdInstanceBody(_a0 []byte) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateVnfdPostBody provides a mock function with given fields: _a0
func (_m *VnfdValidator) ValidateVnfdPostBody(_a0 []byte) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
