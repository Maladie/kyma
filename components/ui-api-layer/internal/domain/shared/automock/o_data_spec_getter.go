// Code generated by mockery v1.0.0. DO NOT EDIT.
package automock

import mock "github.com/stretchr/testify/mock"

import storage "github.com/kyma-project/kyma/components/ui-api-layer/internal/domain/content/storage"

// ODataSpecGetter is an autogenerated mock type for the ODataSpecGetter type
type ODataSpecGetter struct {
	mock.Mock
}

// Find provides a mock function with given fields: kind, id
func (_m *ODataSpecGetter) Find(kind string, id string) (*storage.ODataSpec, error) {
	ret := _m.Called(kind, id)

	var r0 *storage.ODataSpec
	if rf, ok := ret.Get(0).(func(string, string) *storage.ODataSpec); ok {
		r0 = rf(kind, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.ODataSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(kind, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
