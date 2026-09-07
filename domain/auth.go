package domain

import (
	"context"
	"time"
)

type AuthRepository interface {
	RegistrarUsuario(ctx context.Context, email, contraseña, rol, ciudad_origen, nombre_lab string) (string, time.Time, error)
	Login(ctx context.Context, email, contraseña string) (string, time.Time, error)
	Logout(ctx context.Context, token string) error
	Validacion(ctx context.Context, token string) (bool, error)
	CambiarContraseña(ctx context.Context, email, vieja_contraseña, nueva_contraseña string) error
	GetEmail(ctx context.Context, token string) (string, error)
}
