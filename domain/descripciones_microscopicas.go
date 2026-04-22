package domain

type Descripcion_microscopicas struct {
	Descripcion  string		`json:"Descripcion"`
	Diagnositico Diagnostico	`json:"Diagnostico"`
}

type Descripciones_microscopicasRepository interface {
	CreateDescripcionMicroscopica(descripcion_microscopica *Descripcion_microscopicas) error
	GetDescripcionMicroscopicaByProtocolo(protocolo string) (*Descripcion_microscopicas, error)
	// UpdateDescripcionMicroscopica(descripcion_microscopica *Descripcion_microscopicas) error
	DeleteDescripcionMicroscopica(protocolo string) error
}
