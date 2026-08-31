package mocks

import (
	"RATAC/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockUsuarioRepository struct {
	mock.Mock
}

func (m *MockUsuarioRepository) GetByNombre(ctx context.Context, usuario string) (*domain.Usuario, error) {
	args := m.Called(ctx, usuario)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Usuario), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUsuarioRepository) GetByID(ctx context.Context, id int32) (*domain.Usuario, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Usuario), args.Error(1)
	}
	return nil, args.Error(1)
}
