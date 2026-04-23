package domain

import "context"

type Diagnostico struct {
	Descripcion string   `json:"Descripcion"`
	Imagenes    []string `json:"Imagenes"`
}

type DiagnosticoRepository interface {
	// Necesita saber a qué descripción y protocolo asociar el diagnóstico y sus imágenes
	CreateDiagnostico(ctx context.Context, protocolo string, descripcionMicro string, diagnostico *Diagnostico) error

	// Recupera el diagnóstico asociado a una descripción específica
	GetDiagnostico(ctx context.Context, protocolo string, descripcionMicro string) (*Diagnostico, error)

	// Método específico para gestionar la eliminación de imágenes individuales
	DeleteImagen(ctx context.Context, ruta string) error
}
