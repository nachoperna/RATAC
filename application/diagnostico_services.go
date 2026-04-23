package application

import (
	"RATAC/domain"
	"context"
)

type DiagnosticoService struct {
	repo domain.DiagnosticoRepository
}

func NewDiagnosticoService(repo domain.DiagnosticoRepository) *DiagnosticoService {
	return &DiagnosticoService{
		repo: repo,
	}
}

func (s *DiagnosticoService) CreateDiagnostico(ctx context.Context, protocolo string, descripcionMicro string, diagnostico *domain.Diagnostico) error {
	return s.repo.CreateDiagnostico(ctx, protocolo, descripcionMicro, diagnostico)
}

func (s *DiagnosticoService) GetDiagnostico(ctx context.Context, protocolo string, descripcionMicro string) (*domain.Diagnostico, error) {
	return s.repo.GetDiagnostico(ctx, protocolo, descripcionMicro)
}

func (s *DiagnosticoService) DeleteImagen(ctx context.Context, ruta string) error {
	return s.repo.DeleteImagen(ctx, ruta)
}
