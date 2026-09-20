package application

import (
	"RATAC/domain"
	"context"
	"time"
)

type AuthService struct {
	AuthRepo domain.AuthRepository
}

func NewAuthService(AuthRepo domain.AuthRepository) *AuthService {
	return &AuthService{ AuthRepo: AuthRepo, }
}

func (s *AuthService) Registrarse (ctx context.Context, email, contraseña, ciudad_origen, nombre_lab string, matriculas, veterinarios []string) error  {
	return s.AuthRepo.RegistrarSolicitud(ctx, email, contraseña, ciudad_origen, nombre_lab, matriculas, veterinarios)
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
