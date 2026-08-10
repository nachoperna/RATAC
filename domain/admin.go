package domain

import (
	"bytes"
	"context"
)

type AdminRepository interface{
	MapeoDocumento(contenido bytes.Buffer) (*Paciente, error)
	PacienteYaRegistrado(ctx context.Context, protocolo string) bool
}
