// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"
)

// DirectorClient is an autogenerated mock type for the DirectorClient type
type DirectorClient struct {
	mock.Mock
}

// CreateAPIDefinition provides a mock function with given fields: ctx, packageID, apiDefinitionInput
func (_m *DirectorClient) CreateAPIDefinition(ctx context.Context, packageID string, apiDefinitionInput graphql.APIDefinitionInput) (string, error) {
	ret := _m.Called(ctx, packageID, apiDefinitionInput)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.APIDefinitionInput) string); ok {
		r0 = rf(ctx, packageID, apiDefinitionInput)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, graphql.APIDefinitionInput) error); ok {
		r1 = rf(ctx, packageID, apiDefinitionInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateDocument provides a mock function with given fields: ctx, packageID, documentInput
func (_m *DirectorClient) CreateDocument(ctx context.Context, packageID string, documentInput graphql.DocumentInput) (string, error) {
	ret := _m.Called(ctx, packageID, documentInput)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.DocumentInput) string); ok {
		r0 = rf(ctx, packageID, documentInput)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, graphql.DocumentInput) error); ok {
		r1 = rf(ctx, packageID, documentInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateEventDefinition provides a mock function with given fields: ctx, packageID, eventDefinitionInput
func (_m *DirectorClient) CreateEventDefinition(ctx context.Context, packageID string, eventDefinitionInput graphql.EventDefinitionInput) (string, error) {
	ret := _m.Called(ctx, packageID, eventDefinitionInput)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.EventDefinitionInput) string); ok {
		r0 = rf(ctx, packageID, eventDefinitionInput)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, graphql.EventDefinitionInput) error); ok {
		r1 = rf(ctx, packageID, eventDefinitionInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePackage provides a mock function with given fields: ctx, appID, in
func (_m *DirectorClient) CreatePackage(ctx context.Context, appID string, in graphql.PackageCreateInput) (string, error) {
	ret := _m.Called(ctx, appID, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.PackageCreateInput) string); ok {
		r0 = rf(ctx, appID, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, graphql.PackageCreateInput) error); ok {
		r1 = rf(ctx, appID, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAPIDefinition provides a mock function with given fields: ctx, apiID
func (_m *DirectorClient) DeleteAPIDefinition(ctx context.Context, apiID string) error {
	ret := _m.Called(ctx, apiID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, apiID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDocument provides a mock function with given fields: ctx, documentID
func (_m *DirectorClient) DeleteDocument(ctx context.Context, documentID string) error {
	ret := _m.Called(ctx, documentID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, documentID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteEventDefinition provides a mock function with given fields: ctx, eventID
func (_m *DirectorClient) DeleteEventDefinition(ctx context.Context, eventID string) error {
	ret := _m.Called(ctx, eventID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePackage provides a mock function with given fields: ctx, packageID
func (_m *DirectorClient) DeletePackage(ctx context.Context, packageID string) error {
	ret := _m.Called(ctx, packageID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, packageID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetPackage provides a mock function with given fields: ctx, appID, packageID
func (_m *DirectorClient) GetPackage(ctx context.Context, appID string, packageID string) (graphql.PackageExt, error) {
	ret := _m.Called(ctx, appID, packageID)

	var r0 graphql.PackageExt
	if rf, ok := ret.Get(0).(func(context.Context, string, string) graphql.PackageExt); ok {
		r0 = rf(ctx, appID, packageID)
	} else {
		r0 = ret.Get(0).(graphql.PackageExt)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, appID, packageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPackages provides a mock function with given fields: ctx, appID
func (_m *DirectorClient) ListPackages(ctx context.Context, appID string) ([]*graphql.PackageExt, error) {
	ret := _m.Called(ctx, appID)

	var r0 []*graphql.PackageExt
	if rf, ok := ret.Get(0).(func(context.Context, string) []*graphql.PackageExt); ok {
		r0 = rf(ctx, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*graphql.PackageExt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, appID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetApplicationLabel provides a mock function with given fields: ctx, appID, label
func (_m *DirectorClient) SetApplicationLabel(ctx context.Context, appID string, label graphql.LabelInput) error {
	ret := _m.Called(ctx, appID, label)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.LabelInput) error); ok {
		r0 = rf(ctx, appID, label)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePackage provides a mock function with given fields: ctx, packageID, in
func (_m *DirectorClient) UpdatePackage(ctx context.Context, packageID string, in graphql.PackageUpdateInput) error {
	ret := _m.Called(ctx, packageID, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.PackageUpdateInput) error); ok {
		r0 = rf(ctx, packageID, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
