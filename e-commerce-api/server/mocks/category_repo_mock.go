package mocks

import (
	"github.com/stretchr/testify/mock"
)

type CategoryRepositoryMock struct {
	mock.Mock
}

func (m *CategoryRepositoryMock) GetByID(id int) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}
