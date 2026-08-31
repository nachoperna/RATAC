package domain

import "context"

type Rol string

const (
	RolPublico     Rol = "publico"
	RolLaboratorio Rol = "laboratorio"
	RolAdmin       Rol = "admin"
)

type Usuario struct {
	ID           int32  `json:"id"`
	Usuario      string `json:"usuario"`
	PasswordHash string `json:"-"` // nunca se serializa
	Nombre       string `json:"nombre"`
	Rol          Rol    `json:"rol"`
	Activo       bool   `json:"activo"`
}

// UsuarioPublico es el usuario anonimo: sin sesion, sin permisos.
func UsuarioPublico() *Usuario {
	return &Usuario{Rol: RolPublico}
}

// Autenticado distingue a un usuario con sesion del anonimo.
func (u *Usuario) Autenticado() bool {
	return u != nil && u.Rol != RolPublico && u.Rol != ""
}

type UsuarioRepository interface {
	GetByNombre(ctx context.Context, usuario string) (*Usuario, error)
	GetByID(ctx context.Context, id int32) (*Usuario, error)
}
