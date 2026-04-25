package domain

import "context"

type Descripcion_microscopicas struct {
	Descripcion string      `json:"Descripcion"`
	Diagnostico Diagnostico `json:"Diagnostico"`
	TablaGrado	[]Grado_oncologico `json:"Tabla de Grado"`
}

type Descripciones_microscopicasRepository interface {
	// Ahora requiere contexto y el protocolo del paciente por separado
	CreateDescripcionMicroscopica(ctx context.Context, protocolo string, descripcion_microscopica Descripcion_microscopicas) error

	// Devuelve un slice [] porque un paciente puede tener múltiples descripciones
	GetDescripcionMicroscopicaByProtocolo(ctx context.Context, protocolo string) ([]Descripcion_microscopicas, error)

	// Requiere la descripción específica y el protocolo para identificar qué borrar
	DeleteDescripcionMicroscopica(ctx context.Context, descripcion string, protocolo string) error
}
