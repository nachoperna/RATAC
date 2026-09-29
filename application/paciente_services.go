package application

import (
	"RATAC/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lib/pq"
)

// ErrProtocoloFaltante marca un diagnostico sin numero de protocolo. Hay
// documentos que no lo traen adentro y el usuario debe completarlo a mano.
var ErrProtocoloFaltante = errors.New("el diagnostico no tiene numero de protocolo")

// ErrProtocoloDuplicado marca un protocolo que ya existe en la base.
var ErrProtocoloDuplicado = errors.New("el protocolo ya esta registrado")

// esProtocoloDuplicado reconoce la violacion de unicidad que devuelve Postgres
// (SQLSTATE 23505) sobre el protocolo del paciente.
func esProtocoloDuplicado(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	// Si el driver no expone el error tipado, se cae al texto del mensaje
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "pacientes_protocolo_key")
}

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
	// El protocolo identifica al diagnostico y tiene restriccion UNIQUE. Sin esta
	// validacion, un documento sin protocolo se inserta con la cadena vacia y el
	// siguiente falla con un error de clave duplicada que no explica nada.
	if strings.TrimSpace(paciente.Protocolo) == "" {
		return ErrProtocoloFaltante
	}

	err = s.pacienteRepo.InsertarDiagnostico(ctx, paciente)
	if err != nil && esProtocoloDuplicado(err) {
		// El protocolo puede completarse a mano despues de procesar el documento,
		// asi que la verificacion previa no alcanza: el choque aparece recien aca.
		return ErrProtocoloDuplicado
	}
	return err
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

	jsons, err := filepath.Glob(fmt.Sprintf("JSONS/*%s*.json", protocolo))
	if err != nil {
		log.Printf("[DeletePaciente %s] buscar JSON: %v", protocolo, err)
	}

	for _, j := range jsons {
		if err := os.Remove(j); err != nil {
			log.Printf("[DeletePaciente %s] borrar %s: %v", protocolo, j, err)
		}
	}

	imgs, err := filepath.Glob(fmt.Sprintf("IMAGENES/*%s*", protocolo))
	if err != nil {
		log.Printf("[DeletePaciente %s] buscar imagenes: %v", protocolo, err)
	}
	for _, img := range imgs {
		if err := os.Remove(img); err != nil {
			log.Printf("[DeletePaciente %s] borrar %s: %v", protocolo, img, err)
		}
	}
	return nil
}
