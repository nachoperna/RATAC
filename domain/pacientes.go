package domain

import "context"

type Paciente struct {
	Protocolo                   string `json:"Protocolo"`
	Fecha                       string `json:"fecha"`
	Solicitante                 string `json:"Solicitante"`
	Tecnica                     string	`json:"Tecnica"`
	Familia                     string	`json:"Familia"`
	Especie                     string	`json:"Especie"`
	Raza                        string	`json:"Raza"`
	Edad                        string	`json:"Edad"`
	NombrePaciente              string	`json:"Paciente"`
	Antecedentes                string	`json:"Antecedentes"`
	Descripciones_microscopicas []Descripcion_microscopicas `json:"Descripciones microscopicas"`
	DescripcionMacroscopica     string `json:"Descripcion macroscopica"`
	ReferenciasMastocitomas     bool `json:"Referencias mastocitomas"`
}

type PacienteRepository interface {
	CreatePaciente(ctx context.Context, paciente Paciente) error
	// GetPaciente(protocolo string) (Paciente, error)
	// UpdatePaciente(paciente *Paciente) error
	// DeletePaciente(protocolo string) error
	ListPacientes(ctx context.Context) ([]Paciente, error)
}
