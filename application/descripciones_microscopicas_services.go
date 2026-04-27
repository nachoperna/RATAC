package application

import (
	"RATAC/domain"
	"context"
)

type DescripcionMicroscopicaService struct {
	repo domain.Descripcion_microscopicasRepository
}

func NewDescripcionMicroscopicaService(repo domain.Descripcion_microscopicasRepository) *DescripcionMicroscopicaService {
	return &DescripcionMicroscopicaService{
		repo: repo,
	}
}

func (s *DescripcionMicroscopicaService) CreateDescripcionMicroscopica(ctx context.Context, protocolo string, descripcion domain.Descripcion_microscopicas) error {
	return s.repo.CreateDescripcionMicroscopica(ctx, protocolo, descripcion)
}

func (s *DescripcionMicroscopicaService) GetDescripcionesByProtocolo(ctx context.Context, protocolo string) ([]domain.Descripcion_microscopicas, error) {
	return s.repo.GetDescripcionMicroscopicaByProtocolo(ctx, protocolo)
}

func (s *DescripcionMicroscopicaService) DeleteDescripcionMicroscopica(ctx context.Context, descripcion string, protocolo string) error {
	return s.repo.DeleteDescripcionMicroscopica(ctx, descripcion, protocolo)
}

func (s *DescripcionMicroscopicaService) CountDiagnosticos(ctx context.Context) (int64, error) {
	return s.repo.CountDiagnosticos(ctx)
}
