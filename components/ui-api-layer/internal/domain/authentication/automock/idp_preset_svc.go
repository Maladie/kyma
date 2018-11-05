// Code generated by mockery v1.0.0
package automock

import mock "github.com/stretchr/testify/mock"
import pager "github.com/kyma-project/kyma/components/ui-api-layer/internal/pager"
import v1alpha1 "github.com/kyma-project/kyma/components/idppreset/pkg/apis/authentication/v1alpha1"

// idpPresetSvc is an autogenerated mock type for the idpPresetSvc type
type idpPresetSvc struct {
	mock.Mock
}

// Create provides a mock function with given fields: name, issuer, jwksURI
func (_m *idpPresetSvc) Create(name string, issuer string, jwksURI string) (*v1alpha1.IDPPreset, error) {
	ret := _m.Called(name, issuer, jwksURI)

	var r0 *v1alpha1.IDPPreset
	if rf, ok := ret.Get(0).(func(string, string, string) *v1alpha1.IDPPreset); ok {
		r0 = rf(name, issuer, jwksURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.IDPPreset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(name, issuer, jwksURI)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: name
func (_m *idpPresetSvc) Delete(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: name
func (_m *idpPresetSvc) Find(name string) (*v1alpha1.IDPPreset, error) {
	ret := _m.Called(name)

	var r0 *v1alpha1.IDPPreset
	if rf, ok := ret.Get(0).(func(string) *v1alpha1.IDPPreset); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.IDPPreset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: params
func (_m *idpPresetSvc) List(params pager.PagingParams) ([]*v1alpha1.IDPPreset, error) {
	ret := _m.Called(params)

	var r0 []*v1alpha1.IDPPreset
	if rf, ok := ret.Get(0).(func(pager.PagingParams) []*v1alpha1.IDPPreset); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1alpha1.IDPPreset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(pager.PagingParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
