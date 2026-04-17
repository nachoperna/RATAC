package domain

type Paciente struct {
	ID                          int
	Protocolo                   string
	Fecha                       string
	Solicitante                 string
	Tecnica                     string
	Familia                     string
	Especie                     string
	Raza                        string
	Edad                        string
	NombrePaciente              string
	Antecedentes                string
	Descripciones_microscopicas []Descripcion_microscopicas
	DescripcionMacroscopica     string
	ReferenciasMastocitomas     bool
}

type PacienteRepository interface {
	CreatePaciente(paciente *Paciente) error
	GetPacienteByProtocolo(protocolo string) (*Paciente, error)
	// UpdatePaciente(paciente *Paciente) error
	DeletePaciente(protocolo string) error
	ListPacientes() ([]*Paciente, error)
}
