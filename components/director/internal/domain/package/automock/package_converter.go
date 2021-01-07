// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// PackageConverter is an autogenerated mock type for the PackageConverter type
type PackageConverter struct {
	mock.Mock
}

// CreateInputFromGraphQL provides a mock function with given fields: in
func (_m *PackageConverter) CreateInputFromGraphQL(in graphql.PackageCreateInput) (model.PackageCreateInput, error) {
	ret := _m.Called(in)

	var r0 model.PackageCreateInput
	if rf, ok := ret.Get(0).(func(graphql.PackageCreateInput) model.PackageCreateInput); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(model.PackageCreateInput)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(graphql.PackageCreateInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToGraphQL provides a mock function with given fields: in
func (_m *PackageConverter) ToGraphQL(in *model.Package) (*graphql.Package, error) {
	ret := _m.Called(in)

	var r0 *graphql.Package
	if rf, ok := ret.Get(0).(func(*model.Package) *graphql.Package); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.Package)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Package) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateInputFromGraphQL provides a mock function with given fields: in
func (_m *PackageConverter) UpdateInputFromGraphQL(in graphql.PackageUpdateInput) (*model.PackageUpdateInput, error) {
	ret := _m.Called(in)

	var r0 *model.PackageUpdateInput
	if rf, ok := ret.Get(0).(func(graphql.PackageUpdateInput) *model.PackageUpdateInput); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PackageUpdateInput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(graphql.PackageUpdateInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
