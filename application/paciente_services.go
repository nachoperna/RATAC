package application

import (
	"RATAC/domain"
	"context"
)

type PacienteService struct {
	pacienteRepo domain.PacienteRepository
}

func NewPacienteService(pacienteRepo domain.PacienteRepository) *PacienteService {
	return &PacienteService{
		pacienteRepo: pacienteRepo,
	}
}

func (s *PacienteService) ListPacientes(ctx context.Context) ([]domain.Paciente, error) {
	return s.pacienteRepo.ListPacientes(ctx)
}

func (s *PacienteService) ListUltimosPacientes(ctx context.Context) ([]domain.Paciente, []bool, error) {
	return s.pacienteRepo.ListUltimosPacientes(ctx)
}

func (s *PacienteService) CountPacientes(ctx context.Context) (int64, error) {
	return s.pacienteRepo.CountPacientes(ctx)
}

func (s *PacienteService) GetPacienteByNombre(ctx context.Context, nombre string, offset int8) ([]domain.Paciente, int16, error) {
	return s.pacienteRepo.GetPacienteByNombre(ctx, nombre, offset)
}

func (s *PacienteService) GetPacienteByFiltro(ctx context.Context, filtros []domain.Filtro, offset int8) ([]domain.Paciente, int16, error) {
	return s.pacienteRepo.GetPacienteByFiltro(ctx, filtros, offset)
}

func (s *PacienteService) GetAllFromPaciente(ctx context.Context, protocolo string) (*domain.Paciente, error) {
	return s.pacienteRepo.GetAllFromPaciente(ctx, protocolo)
}
