package domain

type Diagnostico struct {
	Descripcion 	string	`json:"Descripcion"`
	Imagenes    	[]string	`json:"Imagenes"`
}

type DiagnosticoRepository interface {
	CreateDiagnostico(Diagnositico *Diagnostico) error
	GetDiagnosticoByProtocolo(protocolo string) (*Diagnostico, error)
	// UpdateDiagnostico(diagnostico *Diagnostico) error
	DeleteDiagnostico(protocolo string) error
}
