package ui_test

import (
	"RATAC/application"
	"RATAC/domain"
	dbrepo "RATAC/infrastructure/DB"
	ui "RATAC/infrastructure/UI"
	"RATAC/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const passwordLogin = "ratac2026"

func adminDePrueba(t *testing.T) *domain.Usuario {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordLogin), bcrypt.MinCost)
	assert.NoError(t, err)
	return &domain.Usuario{
		ID:           1,
		Usuario:      "admin",
		PasswordHash: string(hash),
		Nombre:       "Administrador RATAC",
		Rol:          domain.RolAdmin,
		Activo:       true,
	}
}

func postLogin(form url.Values) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestLogin_credencialesOK_seteaCookieHttpOnlyYLax(t *testing.T) {
	mockRepo := new(mocks.MockUsuarioRepository)
	mockRepo.On("GetByNombre", mock.Anything, "admin").Return(adminDePrueba(t), nil)
	handler := ui.NewAuthHandler(application.NewAuthService(mockRepo))

	rec := httptest.NewRecorder()
	handler.Login(rec, postLogin(url.Values{
		"usuario":  {"admin"},
		"password": {passwordLogin},
		"next":     {"/diagnosticos/alta"},
	}))

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/diagnosticos/alta", rec.Header().Get("Location"))

	cookie := cookieDeRespuesta(rec, ui.NombreCookieSesion)
	assert.NotNil(t, cookie)
	assert.True(t, cookie.HttpOnly, "la cookie debe ser HttpOnly")
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	assert.Equal(t, "/", cookie.Path)
	assert.NotEmpty(t, cookie.Value)
	assert.True(t, cookie.MaxAge > 0)
	// Sin APP_ENV=production debe poder usarse sobre http://localhost.
	assert.False(t, cookie.Secure)
	mockRepo.AssertExpectations(t)
}

func TestLogin_credencialesOKconHTMX_devuelveHXRedirect(t *testing.T) {
	mockRepo := new(mocks.MockUsuarioRepository)
	mockRepo.On("GetByNombre", mock.Anything, "admin").Return(adminDePrueba(t), nil)
	handler := ui.NewAuthHandler(application.NewAuthService(mockRepo))

	req := postLogin(url.Values{"usuario": {"admin"}, "password": {passwordLogin}, "next": {"/admin/panel"}})
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "/admin/panel", rec.Header().Get("HX-Redirect"))
	assert.NotNil(t, cookieDeRespuesta(rec, ui.NombreCookieSesion))
}

// Debe responder 200 y no 401: HTMX por defecto no swapea el body de un 4xx.
func TestLogin_credencialesMalas_devuelve200SinCookie(t *testing.T) {
	mockRepo := new(mocks.MockUsuarioRepository)
	mockRepo.On("GetByNombre", mock.Anything, "admin").Return(adminDePrueba(t), nil)
	handler := ui.NewAuthHandler(application.NewAuthService(mockRepo))

	rec := httptest.NewRecorder()
	handler.Login(rec, postLogin(url.Values{"usuario": {"admin"}, "password": {"incorrecta"}}))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "incorrectos")
	assert.Nil(t, cookieDeRespuesta(rec, ui.NombreCookieSesion), "no debe setear cookie")
}

func TestLogin_usuarioInexistente_devuelve200SinCookie(t *testing.T) {
	mockRepo := new(mocks.MockUsuarioRepository)
	mockRepo.On("GetByNombre", mock.Anything, "fantasma").Return(nil, dbrepo.ErrUsuarioNoEncontrado)
	handler := ui.NewAuthHandler(application.NewAuthService(mockRepo))

	rec := httptest.NewRecorder()
	handler.Login(rec, postLogin(url.Values{"usuario": {"fantasma"}, "password": {passwordLogin}}))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "incorrectos")
	assert.Nil(t, cookieDeRespuesta(rec, ui.NombreCookieSesion))
}

// Open redirect: un next externo debe caer al destino por defecto.
func TestLogin_nextMalicioso_caeADestinoPorDefecto(t *testing.T) {
	casos := []string{"//evil.com", "https://evil.com", "http://evil.com/x", "evil.com", "\\\\evil.com"}

	for _, next := range casos {
		mockRepo := new(mocks.MockUsuarioRepository)
		mockRepo.On("GetByNombre", mock.Anything, "admin").Return(adminDePrueba(t), nil)
		handler := ui.NewAuthHandler(application.NewAuthService(mockRepo))

		rec := httptest.NewRecorder()
		handler.Login(rec, postLogin(url.Values{
			"usuario":  {"admin"},
			"password": {passwordLogin},
			"next":     {next},
		}))

		assert.Equal(t, "/admin/panel", rec.Header().Get("Location"), "next %q", next)
	}
}

func TestLogin_nextInterno_seRespeta(t *testing.T) {
	mockRepo := new(mocks.MockUsuarioRepository)
	mockRepo.On("GetByNombre", mock.Anything, "admin").Return(adminDePrueba(t), nil)
	handler := ui.NewAuthHandler(application.NewAuthService(mockRepo))

	rec := httptest.NewRecorder()
	handler.Login(rec, postLogin(url.Values{
		"usuario":  {"admin"},
		"password": {passwordLogin},
		"next":     {"/paciente/protocolo/CAN-001"},
	}))

	assert.Equal(t, "/paciente/protocolo/CAN-001", rec.Header().Get("Location"))
}

func TestLogout_borraCookie(t *testing.T) {
	handler := ui.NewAuthHandler(application.NewAuthService(new(mocks.MockUsuarioRepository)))

	rec := httptest.NewRecorder()
	handler.Logout(rec, httptest.NewRequest(http.MethodPost, "/logout", nil))

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))

	cookie := cookieDeRespuesta(rec, ui.NombreCookieSesion)
	assert.NotNil(t, cookie)
	assert.True(t, cookie.MaxAge < 0, "MaxAge negativo borra la cookie")
	assert.Empty(t, cookie.Value)
}

func TestShowLogin_yaAutenticado_redirige(t *testing.T) {
	handler := ui.NewAuthHandler(application.NewAuthService(new(mocks.MockUsuarioRepository)))

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req = req.WithContext(domain.ConUsuarioCtx(req.Context(),
		&domain.Usuario{ID: 1, Usuario: "admin", Rol: domain.RolAdmin}))
	rec := httptest.NewRecorder()

	handler.ShowLogin(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/admin/panel", rec.Header().Get("Location"))
}

func TestShowLogin_anonimo_renderizaFormulario(t *testing.T) {
	handler := ui.NewAuthHandler(application.NewAuthService(new(mocks.MockUsuarioRepository)))

	rec := httptest.NewRecorder()
	handler.ShowLogin(rec, httptest.NewRequest(http.MethodGet, "/login?next=%2Fadmin%2Fpanel", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `id="login-form"`)
	assert.Contains(t, body, `name="usuario"`)
	assert.Contains(t, body, `name="password"`)
	assert.Contains(t, body, "/admin/panel")
}
