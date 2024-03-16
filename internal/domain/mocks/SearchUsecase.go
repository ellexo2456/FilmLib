// Code generated by mockery v2.35.2. DO NOT EDIT.

package mocks

import (
	domain "2023_2_Holi/domain"

	mock "github.com/stretchr/testify/mock"
)

// SearchUsecase is an autogenerated mock type for the SearchUsecase type
type SearchUsecase struct {
	mock.Mock
}

// GetSearchData provides a mock function with given fields: searchStr
func (_m *SearchUsecase) GetSearchData(searchStr string) (domain.SearchData, error) {
	ret := _m.Called(searchStr)

	var r0 domain.SearchData
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.SearchData, error)); ok {
		return rf(searchStr)
	}
	if rf, ok := ret.Get(0).(func(string) domain.SearchData); ok {
		r0 = rf(searchStr)
	} else {
		r0 = ret.Get(0).(domain.SearchData)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(searchStr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSearchUsecase creates a new instance of SearchUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSearchUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *SearchUsecase {
	mock := &SearchUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}