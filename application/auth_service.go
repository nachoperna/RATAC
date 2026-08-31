package application

import (
	"RATAC/domain"
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	DuracionSesion = 8 * time.Hour
	emisorToken    = "ratac"
	metodoFirma    = "HS256"
)

// ErrCredencialesInvalidas es el unico error que Login devuelve ante un fallo
// de autenticacion, sea por usuario inexistente o por password incorrecta:
// distinguirlos filtraria que usuarios existen.
var ErrCredencialesInvalidas = errors.New("usuario o contraseña incorrectos")

// hashDummy es un bcrypt valido contra el que se compara cuando el usuario no
// existe, para que el tiempo de respuesta no delate su ausencia.
const hashDummy = "$2a$10$wG9/z1UV22rUkYtGJ.Qfvu8fBA32G3zyd3WM0Y0G4RKE/AtAqAICW"

type AuthService struct {
	repo   domain.UsuarioRepository
	secret []byte
}

func NewAuthService(repo domain.UsuarioRepository) *AuthService {
	return &AuthService{
		repo:   repo,
		secret: []byte(getJWTSecret()),
	}
}

func getJWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "ratac-secret-de-desarrollo-cambiar-en-produccion"
}

type ClaimsRATAC struct {
	Usuario string `json:"usr"`
	Nombre  string `json:"nom"`
	Rol     string `json:"rol"`
	jwt.RegisteredClaims
}

func (s *AuthService) Login(ctx context.Context, usuario, password string) (*domain.Usuario, error) {
	u, err := s.repo.GetByNombre(ctx, usuario)
	if err != nil || u == nil {
		// Comparacion contra un hash dummy para no filtrar existencia por timing.
		_ = bcrypt.CompareHashAndPassword([]byte(hashDummy), []byte(password))
		return nil, ErrCredencialesInvalidas
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrCredencialesInvalidas
	}
	return u, nil
}

func (s *AuthService) EmitirToken(u *domain.Usuario) (string, error) {
	if u == nil {
		return "", errors.New("no se puede emitir un token sin usuario")
	}
	ahora := time.Now()
	claims := ClaimsRATAC{
		Usuario: u.Usuario,
		Nombre:  u.Nombre,
		Rol:     string(u.Rol),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(u.ID)),
			Issuer:    emisorToken,
			IssuedAt:  jwt.NewNumericDate(ahora),
			ExpiresAt: jwt.NewNumericDate(ahora.Add(DuracionSesion)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidarToken reconstruye el usuario desde los claims sin ir a la BD (una
// query menos por request). Costo: revocar o cambiar un rol tarda hasta
// DuracionSesion en propagarse; cambiar a repo.GetByID aca es una linea si
// mas adelante hace falta revocacion inmediata.
func (s *AuthService) ValidarToken(tokenStr string) (*domain.Usuario, error) {
	claims := &ClaimsRATAC{}
	// WithValidMethods es obligatorio: sin el, un token con alg "none" pasa.
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return s.secret, nil
	},
		jwt.WithValidMethods([]string{metodoFirma}),
		jwt.WithIssuer(emisorToken),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	rol := domain.Rol(claims.Rol)
	if _, ok := domain.PermisosPorRol[rol]; !ok || rol == domain.RolPublico {
		return nil, errors.New("rol inválido en el token")
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, errors.New("subject inválido en el token")
	}

	return &domain.Usuario{
		ID:      int32(id),
		Usuario: claims.Usuario,
		Nombre:  claims.Nombre,
		Rol:     rol,
		Activo:  true,
	}, nil
}
