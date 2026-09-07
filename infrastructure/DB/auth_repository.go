package db

import (
	sqlc "RATAC/DB/sqlc"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	queries *sqlc.Queries
	db *sql.DB
}

func NewAuthRepository(queries *sqlc.Queries, db *sql.DB) *AuthRepository {
	return &AuthRepository{
		queries: queries,
		db: db,
	}
}

func (r *AuthRepository) RegistrarUsuario(ctx context.Context, email, contraseña, rol, ciudad_origen, nombre_lab string) (string, time.Time, error)  {
	hash, err := bcrypt.GenerateFromPassword([]byte(contraseña), bcrypt.DefaultCost)
	if err != nil {
		return "", time.Time{}, err
	}

	usuario, err := r.queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
		Email: email,
		ContraseñaHash: string(hash),
		Rol: sqlc.Roles(rol),
		CiudadOrigen: ciudad_origen,
		NombreLab: nombre_lab,
	})
	if err != nil {
		return "", time.Time{}, err
	}

	token, expiracion, err := r.NuevaSesion(ctx, usuario.ID)
	return token, expiracion, err
}

func (r *AuthRepository) Login(ctx context.Context, email, contraseña string) (string, time.Time, error)  {
	usuario, err := r.queries.Login(ctx, email)
	if err != nil {
		return "", time.Time{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(usuario.ContraseñaHash), []byte(contraseña))
	if err != nil {
		return "", time.Time{}, err
	}

	token, expiracion, err := r.NuevaSesion(ctx, usuario.ID)
	return token, expiracion, err
}
	
func (r *AuthRepository) Validacion(ctx context.Context, token string) (bool, error)  {
	sesion_activa, err := r.queries.SesionActiva(ctx, token)
	if err != nil || !sesion_activa {
		return false, err
	}
	return true, nil
}

func (r *AuthRepository) NuevaSesion(ctx context.Context, id uuid.UUID) (string, time.Time, error) {
	token := uuid.New().String()
	sesion, err := r.queries.NuevaSesion(ctx, sqlc.NuevaSesionParams{
		Token: token,
		IDUsuario: id,
	})
	return token, sesion.Expiracion, err
}

func (r *AuthRepository) Logout(ctx context.Context, token string) error  {
	return r.queries.DeleteSesion(ctx, token)
}

func (r *AuthRepository) CambiarContraseña(ctx context.Context, email, vieja_contraseña, nueva_contraseña string) error {
	usuario, err := r.queries.GetUsuario(ctx, email)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(usuario.ContraseñaHash), []byte(vieja_contraseña))
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(nueva_contraseña), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = r.queries.UpdateContraseña(ctx, sqlc.UpdateContraseñaParams{
		Email: email,
		ContraseñaHash: string(hash),
	})
	return err
}

func (r *AuthRepository) GetEmail(ctx context.Context, token string) (string, error) {
	return r.queries.GetEmail(ctx, token)
}
