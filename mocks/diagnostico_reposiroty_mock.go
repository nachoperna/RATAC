package mocks

import (
	"RATAC/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockDiagnosticoRepository simula la interfaz domain.DiagnosticoRepository
type MockDiagnosticoRepository struct {
	mock.Mock
}

func (m *MockDiagnosticoRepository) CreateDiagnostico(ctx context.Context, protocolo string, descripcionMicro string, diagnostico domain.Diagnostico) error {
	args := m.Called(ctx, protocolo, descripcionMicro, diagnostico)
	return args.Error(0)
}

func (m *MockDiagnosticoRepository) GetDiagnostico(ctx context.Context, protocolo string, descripcionMicro string) (*domain.Diagnostico, error) {
	args := m.Called(ctx, protocolo, descripcionMicro)
	// Como retornamos un puntero, validamos si es nil
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Diagnostico), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDiagnosticoRepository) DeleteImagen(ctx context.Context, ruta string) error {
	args := m.Called(ctx, ruta)
	return args.Error(0)
}

func (m *MockDiagnosticoRepository) CountImagenes(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
