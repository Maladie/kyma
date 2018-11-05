// Code generated by mockery v1.0.0
package automock

import mock "github.com/stretchr/testify/mock"

import v1alpha1 "github.com/kyma-project/kyma/components/remote-environment-broker/pkg/apis/applicationconnector/v1alpha1"

// eventActivationLister is an autogenerated mock type for the eventActivationLister type
type eventActivationLister struct {
	mock.Mock
}

// List provides a mock function with given fields: environment
func (_m *eventActivationLister) List(environment string) ([]*v1alpha1.EventActivation, error) {
	ret := _m.Called(environment)

	var r0 []*v1alpha1.EventActivation
	if rf, ok := ret.Get(0).(func(string) []*v1alpha1.EventActivation); ok {
		r0 = rf(environment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1alpha1.EventActivation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(environment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
