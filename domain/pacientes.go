package domain

import "context"

type Paciente struct {
	Protocolo                   string                      `json:"Protocolo"`
	Fecha                       string                      `json:"fecha"`
	Solicitante                 string                      `json:"Solicitante"`
	Tecnica                     string                      `json:"Técnica"`
	Familia                     string                      `json:"Familia"`
	Especie                     string                      `json:"Especie"`
	Raza                        string                      `json:"Raza"`
	Edad                        string                      `json:"Edad"`
	NombrePaciente              string                      `json:"Paciente"`
	ReferenciasMastocitomas     bool                        `json:"Referencias mastocitomas"`
	Antecedentes                string                      `json:"Material remitido - Antecedentes"`
	DescripcionMacroscopica     string                      `json:"Descripción macroscópica"`
	Descripciones_microscopicas []Descripcion_microscopicas `json:"Descripción microscópica"`
}

type Filtro struct {
	Logica   string `json:"logica"`
	Campo    string `json:"campo"`
	Operador string `json:"operador"`
	Valores  []string `json:"valores"`
	Not      bool   `json:"not"`
	Multiple bool   `json:"multiple"`
}

// Lista blanca de columnas permitidas
var ColumnasPermitidas = map[string]bool{
	"Especie": true,
	"Raza":    true,
	"Edad":    true,
	"Paciente": true,
	"Protocolo": true,
}

type PacienteRepository interface {
	CreatePaciente(ctx context.Context, paciente Paciente) error
	// GetPaciente(protocolo string) (Paciente, error)
	// UpdatePaciente(paciente *Paciente) error
	// DeletePaciente(protocolo string) error
	ListUltimosPacientes(ctx context.Context) ([]Paciente, error)
	ListPacientes(ctx context.Context) ([]Paciente, error)
	CountPacientes(ctx context.Context) (int64, error)
	GetPacienteByNombre(ctx context.Context, nombre string) ([]Paciente, int16, error)
	GetPacienteByFiltro(ctx context.Context, filtros []Filtro) ([]Paciente, int16, error)
}
