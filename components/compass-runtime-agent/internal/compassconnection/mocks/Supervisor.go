// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	v1alpha1 "github.com/kyma-project/kyma/components/compass-runtime-agent/pkg/apis/compass/v1alpha1"
	mock "github.com/stretchr/testify/mock"
)

// Supervisor is an autogenerated mock type for the Supervisor type
type Supervisor struct {
	mock.Mock
}

// InitializeCompassConnection provides a mock function with given fields:
func (_m *Supervisor) InitializeCompassConnection() (*v1alpha1.CompassConnection, error) {
	ret := _m.Called()

	var r0 *v1alpha1.CompassConnection
	if rf, ok := ret.Get(0).(func() *v1alpha1.CompassConnection); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.CompassConnection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SynchronizeWithCompass provides a mock function with given fields: connection
func (_m *Supervisor) SynchronizeWithCompass(connection *v1alpha1.CompassConnection) (*v1alpha1.CompassConnection, error) {
	ret := _m.Called(connection)

	var r0 *v1alpha1.CompassConnection
	if rf, ok := ret.Get(0).(func(*v1alpha1.CompassConnection) *v1alpha1.CompassConnection); ok {
		r0 = rf(connection)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.CompassConnection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1alpha1.CompassConnection) error); ok {
		r1 = rf(connection)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
