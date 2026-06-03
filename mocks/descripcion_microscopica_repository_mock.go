package mocks

import (
	"RATAC/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockDescripcionMicroscopicaRepository struct {
	mock.Mock
}

func (m *MockDescripcionMicroscopicaRepository) CreateDescripcionMicroscopica(ctx context.Context, protocolo string, descripcion domain.Descripcion_microscopicas) error {
	args := m.Called(ctx, protocolo, descripcion)
	return args.Error(0)
}

func (m *MockDescripcionMicroscopicaRepository) GetDescripcionMicroscopicaByProtocolo(ctx context.Context, protocolo string) ([]domain.Descripcion_microscopicas, error) {
	args := m.Called(ctx, protocolo)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Descripcion_microscopicas), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDescripcionMicroscopicaRepository) DeleteDescripcionMicroscopica(ctx context.Context, descripcion string, protocolo string) error {
	args := m.Called(ctx, descripcion, protocolo)
	return args.Error(0)
}

func (m *MockDescripcionMicroscopicaRepository) CountDiagnosticos(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
