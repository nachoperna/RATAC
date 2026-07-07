package domain

import "bytes"

type AdminRepository interface{
	MapeoDocumento(contenido bytes.Buffer) (*Paciente, error)
}
