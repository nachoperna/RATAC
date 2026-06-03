package mocks

import (
	"RATAC/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockPacienteRepository struct {
	mock.Mock
}

func (m *MockPacienteRepository) InsertarDiagnostico(ctx context.Context, paciente domain.Paciente) error {
	args := m.Called(ctx, paciente)
	return args.Error(0)
}

func (m *MockPacienteRepository) ListUltimosPacientes(ctx context.Context) ([]domain.Paciente, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Paciente), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPacienteRepository) ListPacientes(ctx context.Context) ([]domain.Paciente, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Paciente), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPacienteRepository) CountPacientes(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPacienteRepository) GetPacienteByNombre(ctx context.Context, nombre string, offset int8) ([]domain.Paciente, int16, error) {
	args := m.Called(ctx, nombre, offset)
	var pacientes []domain.Paciente
	if args.Get(0) != nil {
		pacientes = args.Get(0).([]domain.Paciente)
	}
	return pacientes, args.Get(1).(int16), args.Error(2)
}

func (m *MockPacienteRepository) GetPacienteByFiltro(ctx context.Context, filtros []domain.Filtro, offset int8) ([]domain.Paciente, int16, error) {
	args := m.Called(ctx, filtros, offset)
	var pacientes []domain.Paciente
	if args.Get(0) != nil {
		pacientes = args.Get(0).([]domain.Paciente)
	}
	return pacientes, args.Get(1).(int16), args.Error(2)
}

func (m *MockPacienteRepository) GetAllFromPaciente(ctx context.Context, protocolo string) (*domain.Paciente, error) {
	args := m.Called(ctx, protocolo)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Paciente), args.Error(1)
	}
	return nil, args.Error(1)
}
