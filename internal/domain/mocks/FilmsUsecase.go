// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	domain "github.com/ellexo2456/FilmLib/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// FilmsUsecase is an autogenerated mock type for the FilmsUsecase type
type FilmsUsecase struct {
	mock.Mock
}

// Add provides a mock function with given fields: film
func (_m *FilmsUsecase) Add(film domain.Film) (int, error) {
	ret := _m.Called(film)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.Film) (int, error)); ok {
		return rf(film)
	}
	if rf, ok := ret.Get(0).(func(domain.Film) int); ok {
		r0 = rf(film)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(domain.Film) error); ok {
		r1 = rf(film)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: title, releaseDate
func (_m *FilmsUsecase) GetAll(title domain.SortDirection, releaseDate domain.SortDirection) ([]domain.Film, error) {
	ret := _m.Called(title, releaseDate)

	var r0 []domain.Film
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.SortDirection, domain.SortDirection) ([]domain.Film, error)); ok {
		return rf(title, releaseDate)
	}
	if rf, ok := ret.Get(0).(func(domain.SortDirection, domain.SortDirection) []domain.Film); ok {
		r0 = rf(title, releaseDate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Film)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.SortDirection, domain.SortDirection) error); ok {
		r1 = rf(title, releaseDate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Modify provides a mock function with given fields: film
func (_m *FilmsUsecase) Modify(film domain.Film) (domain.Film, error) {
	ret := _m.Called(film)

	var r0 domain.Film
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.Film) (domain.Film, error)); ok {
		return rf(film)
	}
	if rf, ok := ret.Get(0).(func(domain.Film) domain.Film); ok {
		r0 = rf(film)
	} else {
		r0 = ret.Get(0).(domain.Film)
	}

	if rf, ok := ret.Get(1).(func(domain.Film) error); ok {
		r1 = rf(film)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: id
func (_m *FilmsUsecase) Remove(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Search provides a mock function with given fields: searchStr
func (_m *FilmsUsecase) Search(searchStr string) ([]domain.Film, error) {
	ret := _m.Called(searchStr)

	var r0 []domain.Film
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.Film, error)); ok {
		return rf(searchStr)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.Film); ok {
		r0 = rf(searchStr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Film)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(searchStr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFilmsUsecase creates a new instance of FilmsUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFilmsUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *FilmsUsecase {
	mock := &FilmsUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
