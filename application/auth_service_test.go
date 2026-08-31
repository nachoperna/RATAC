package application_test

import (
	"RATAC/application"
	"RATAC/domain"
	dbrepo "RATAC/infrastructure/DB"
	"RATAC/mocks"
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const passwordPrueba = "ratac2026"

func usuarioDePrueba(t *testing.T, rol domain.Rol) *domain.Usuario {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordPrueba), bcrypt.MinCost)
	assert.NoError(t, err)
	return &domain.Usuario{
		ID:           1,
		Usuario:      "admin",
		PasswordHash: string(hash),
		Nombre:       "Administrador RATAC",
		Rol:          rol,
		Activo:       true,
	}
}

func TestLogin_passwordCorrecta_devuelveUsuario(t *testing.T) {
	mockRepo := new(mocks.MockUsuarioRepository)
	service := application.NewAuthService(mockRepo)
	esperado := usuarioDePrueba(t, domain.RolAdmin)

	mockRepo.On("GetByNombre", mock.Anything, "admin").Return(esperado, nil)

	u, err := service.Login(context.Background(), "admin", passwordPrueba)

	assert.NoError(t, err)
	assert.Equal(t, esperado, u)
	mockRepo.AssertExpectations(t)
}

func TestLogin_passwordIncorrecta_devuelveErrCredenciales(t *testing.T) {
	mockRepo := new(mocks.MockUsuarioRepository)
	service := application.NewAuthService(mockRepo)

	mockRepo.On("GetByNombre", mock.Anything, "admin").Return(usuarioDePrueba(t, domain.RolAdmin), nil)

	u, err := service.Login(context.Background(), "admin", "incorrecta")

	assert.Nil(t, u)
	assert.ErrorIs(t, err, application.ErrCredencialesInvalidas)
}

// Usuario inexistente y password mala deben devolver el mismo error, para no
// filtrar que usuarios existen.
func TestLogin_usuarioInexistente_devuelveMismoErrorQuePasswordMala(t *testing.T) {
	mockRepo := new(mocks.MockUsuarioRepository)
	service := application.NewAuthService(mockRepo)

	mockRepo.On("GetByNombre", mock.Anything, "fantasma").Return(nil, dbrepo.ErrUsuarioNoEncontrado)

	u, err := service.Login(context.Background(), "fantasma", passwordPrueba)

	assert.Nil(t, u)
	assert.ErrorIs(t, err, application.ErrCredencialesInvalidas)
	assert.Equal(t, application.ErrCredencialesInvalidas.Error(), err.Error())
}

func TestEmitirYValidarToken_preservaRol(t *testing.T) {
	service := application.NewAuthService(new(mocks.MockUsuarioRepository))

	for _, rol := range []domain.Rol{domain.RolAdmin, domain.RolLaboratorio} {
		original := usuarioDePrueba(t, rol)

		token, err := service.EmitirToken(original)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		u, err := service.ValidarToken(token)
		assert.NoError(t, err)
		assert.Equal(t, rol, u.Rol)
		assert.Equal(t, original.ID, u.ID)
		assert.Equal(t, original.Usuario, u.Usuario)
		assert.Equal(t, original.Nombre, u.Nombre)
		// El hash nunca viaja en el token.
		assert.Empty(t, u.PasswordHash)
	}
}

func TestValidarToken_otroSecret_devuelveError(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret-del-emisor")
	emisor := application.NewAuthService(new(mocks.MockUsuarioRepository))
	token, err := emisor.EmitirToken(usuarioDePrueba(t, domain.RolAdmin))
	assert.NoError(t, err)

	t.Setenv("JWT_SECRET", "un-secret-totalmente-distinto")
	validador := application.NewAuthService(new(mocks.MockUsuarioRepository))

	u, err := validador.ValidarToken(token)

	assert.Nil(t, u)
	assert.Error(t, err)
}

func TestValidarToken_expirado_devuelveError(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret-de-prueba")
	service := application.NewAuthService(new(mocks.MockUsuarioRepository))

	// Token firmado con el mismo secret pero ya vencido.
	pasado := time.Now().Add(-2 * time.Hour)
	claims := application.ClaimsRATAC{
		Usuario: "admin",
		Nombre:  "Administrador RATAC",
		Rol:     string(domain.RolAdmin),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(1),
			Issuer:    "ratac",
			IssuedAt:  jwt.NewNumericDate(pasado),
			ExpiresAt: jwt.NewNumericDate(pasado.Add(time.Hour)),
		},
	}
	firmado, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret-de-prueba"))
	assert.NoError(t, err)

	u, err := service.ValidarToken(firmado)

	assert.Nil(t, u)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

// Sin jwt.WithValidMethods un token con alg "none" pasaria la validacion.
func TestValidarToken_algNone_devuelveError(t *testing.T) {
	service := application.NewAuthService(new(mocks.MockUsuarioRepository))

	claims := application.ClaimsRATAC{
		Usuario: "atacante",
		Rol:     string(domain.RolAdmin),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			Issuer:    "ratac",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	firmado, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	u, err := service.ValidarToken(firmado)

	assert.Nil(t, u)
	assert.Error(t, err)
}

func TestValidarToken_otroIssuer_devuelveError(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret-de-prueba")
	service := application.NewAuthService(new(mocks.MockUsuarioRepository))

	claims := application.ClaimsRATAC{
		Usuario: "admin",
		Rol:     string(domain.RolAdmin),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			Issuer:    "otro-sistema",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	firmado, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret-de-prueba"))
	assert.NoError(t, err)

	u, err := service.ValidarToken(firmado)

	assert.Nil(t, u)
	assert.Error(t, err)
}

// Un token con rol "publico" o con un rol desconocido no debe producir sesion.
func TestValidarToken_rolInvalido_devuelveError(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret-de-prueba")
	service := application.NewAuthService(new(mocks.MockUsuarioRepository))

	for _, rol := range []string{string(domain.RolPublico), "superusuario", ""} {
		claims := application.ClaimsRATAC{
			Usuario: "atacante",
			Rol:     rol,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "1",
				Issuer:    "ratac",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}
		firmado, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret-de-prueba"))
		assert.NoError(t, err)

		u, err := service.ValidarToken(firmado)
		assert.Nil(t, u, "rol %q", rol)
		assert.Error(t, err, "rol %q", rol)
	}
}

func TestValidarToken_basura_devuelveError(t *testing.T) {
	service := application.NewAuthService(new(mocks.MockUsuarioRepository))
	u, err := service.ValidarToken("esto-no-es-un-jwt")
	assert.Nil(t, u)
	assert.Error(t, err)
}
