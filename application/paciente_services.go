package application

import (
	"RATAC/domain"
	"context"
)

type PacienteService struct{
	pacienteRepo domain.PacienteRepository
}

func NewPacienteService(pacienteRepo domain.PacienteRepository) *PacienteService{
	return &PacienteService{
		pacienteRepo: pacienteRepo,
	}
}

func (s *PacienteService) CreatePaciente(ctx context.Context, paciente domain.Paciente) error {
	return s.pacienteRepo.CreatePaciente(ctx, paciente)
}

func (s *PacienteService) ListPacientes(ctx context.Context) ([]domain.Paciente, error) {
	return s.pacienteRepo.ListPacientes(ctx)
}
