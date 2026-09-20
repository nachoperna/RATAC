package domain

import (
	"bytes"
	"context"
)

type Solicitud struct {
	Email string
	Nombre string 
	Ciudad string 
	Veterinarios []Veterinarios
}

type Veterinarios struct {
	Matricula int32
	Nombre string
}

type AdminRepository interface{
	MapeoDocumento(contenido bytes.Buffer) (*Paciente, error)
	PacienteYaRegistrado(ctx context.Context, protocolo string) bool
	GetUltimosDiagnosticosCargados(ctx context.Context, offset int8, token string) ([]Paciente, int16, error)
	GetSolicitudes(ctx context.Context) ([]Solicitud, error)
}
