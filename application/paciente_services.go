package application

import (
	"RATAC/domain"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type PacienteService struct {
	pacienteRepo domain.PacienteRepository
}

func NewPacienteService(pacienteRepo domain.PacienteRepository) *PacienteService {
	return &PacienteService{
		pacienteRepo: pacienteRepo,
	}
}
func (s *PacienteService) InsertarDiagnostico(ctx context.Context, nombre_archivo string) error {
	archivo, err := os.Open(fmt.Sprintf("JSONS/%s.json", nombre_archivo))
	if err != nil {
		return err
	}
	var paciente domain.Paciente
	err = json.NewDecoder(archivo).Decode(&paciente)
	if err != nil {
		return err
	}
	return s.pacienteRepo.InsertarDiagnostico(ctx, paciente)
}

func (s *PacienteService) ListPacientes(ctx context.Context, offset int8) ([]domain.Paciente, int16, error) {
	return s.pacienteRepo.ListPacientes(ctx, offset)
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

func (s *PacienteService) DeletePaciente(ctx context.Context, protocolo string) error {
	err := s.pacienteRepo.DeletePaciente(ctx, protocolo)
	if err != nil {
		return err
	}
	json, err := filepath.Glob(fmt.Sprintf("JSONS/*%s*.json", protocolo))
	err = os.Remove(json[0])
	imgs, err := filepath.Glob(fmt.Sprintf("IMAGENES/*%s*", protocolo))
	for _, img := range imgs {
		err = os.Remove(img)
	}
	if err != nil {
		return err
	}
	return nil
}
