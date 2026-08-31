package domain

import (
	"bytes"
	"context"
)

type AdminRepository interface{
	MapeoDocumento(contenido bytes.Buffer) (*Paciente, error)
	PacienteYaRegistrado(ctx context.Context, protocolo string) bool
	GetUltimosDiagnosticosCargados(ctx context.Context, offset int8) ([]Paciente, int16, error)
}
