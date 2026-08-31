package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"context"
	"database/sql"
	"errors"
)

type UsuarioRepository struct {
	queries *sqlc.Queries
}

func NewUsuarioRepository(queries *sqlc.Queries) *UsuarioRepository {
	return &UsuarioRepository{queries: queries}
}

// ErrUsuarioNoEncontrado se devuelve ante un usuario inexistente o inactivo.
// El servicio de auth lo traduce al mismo error que una password incorrecta
// para no filtrar que usuarios existen.
var ErrUsuarioNoEncontrado = errors.New("usuario no encontrado")

func (r *UsuarioRepository) GetByNombre(ctx context.Context, usuario string) (*domain.Usuario, error) {
	u, err := r.queries.GetUsuarioByNombre(ctx, usuario)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUsuarioNoEncontrado
		}
		return nil, err
	}
	return usuarioORM(u), nil
}

func (r *UsuarioRepository) GetByID(ctx context.Context, id int32) (*domain.Usuario, error) {
	u, err := r.queries.GetUsuarioByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUsuarioNoEncontrado
		}
		return nil, err
	}
	return usuarioORM(u), nil
}

func usuarioORM(u sqlc.Usuario) *domain.Usuario {
	return &domain.Usuario{
		ID:           u.ID,
		Usuario:      u.Usuario,
		PasswordHash: u.PasswordHash,
		Nombre:       u.Nombre,
		Rol:          domain.Rol(u.Rol),
		Activo:       u.Activo,
	}
}
