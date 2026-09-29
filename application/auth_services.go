package application

import (
	"RATAC/domain"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
)

// ErrEmailYaRegistrado marca un alta con un correo que ya tiene cuenta. Es el
// unico fallo del registro que el usuario puede resolver por su cuenta.
var ErrEmailYaRegistrado = errors.New("el correo ya tiene una cuenta registrada")

// esEmailDuplicado reconoce la violacion de unicidad de Postgres (SQLSTATE 23505)
// sobre usuarios.email.
func esEmailDuplicado(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "usuarios_email_key")
}

type AuthService struct {
	AuthRepo domain.AuthRepository
}

func NewAuthService(AuthRepo domain.AuthRepository) *AuthService {
	return &AuthService{ AuthRepo: AuthRepo, }
}

func (s *AuthService) Registrarse (ctx context.Context, email, contraseña, ciudad_origen, nombre_lab string, matriculas, veterinarios []string) error  {
	err := s.AuthRepo.RegistrarSolicitud(ctx, email, contraseña, ciudad_origen, nombre_lab, matriculas, veterinarios)
	if err != nil && esEmailDuplicado(err) {
		return ErrEmailYaRegistrado
	}
	return err
}
func (s *AuthService) Login (ctx context.Context, email, contraseña string) (string, time.Time, error)  {
	return s.AuthRepo.Login(ctx, email, contraseña)
}
func (s *AuthService) Logout (ctx context.Context, token string) error  {
	return s.AuthRepo.Logout(ctx, token)
}
func (s *AuthService) Validacion (ctx context.Context, token string) (bool, error)  {
	return s.AuthRepo.Validacion(ctx, token)
}
func (s *AuthService) CambiarContraseña (ctx context.Context, email, vieja_contraseña, nueva_contraseña string) error  {
	return s.AuthRepo.CambiarContraseña(ctx, email, vieja_contraseña, nueva_contraseña)
}
func (s *AuthService) GetEmail (ctx context.Context, token string) (string, error)  {
	return s.AuthRepo.GetEmail(ctx, token)
}
