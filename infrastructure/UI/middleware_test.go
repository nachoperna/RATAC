package ui_test

import (
	"RATAC/application"
	"RATAC/domain"
	ui "RATAC/infrastructure/UI"
	"RATAC/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func handlerCentinela(invocado *bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		*invocado = true
		w.WriteHeader(http.StatusOK)
	}
}

// requestConUsuario simula lo que hace ConUsuario, sin pasar por la cookie.
func requestConUsuario(metodo, ruta string, u *domain.Usuario) *http.Request {
	req := httptest.NewRequest(metodo, ruta, nil)
	if u != nil {
		req = req.WithContext(domain.ConUsuarioCtx(req.Context(), u))
	}
	return req
}

func TestRequiere_anonimo_redirigeALoginYNoInvocaHandler(t *testing.T) {
	invocado := false
	protegido := ui.Requiere(domain.PermVerPanel, handlerCentinela(&invocado))

	req := httptest.NewRequest(http.MethodGet, "/admin/panel", nil)
	rec := httptest.NewRecorder()

	protegido(rec, req)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/login?next=%2Fadmin%2Fpanel", rec.Header().Get("Location"))
	assert.False(t, invocado, "el handler protegido no debe ejecutarse")
}

// Ante HTMX hay que responder 200 + HX-Redirect: un 302 haria que HTMX
// inyecte el HTML del login dentro del target del swap.
func TestRequiere_anonimoConHTMX_devuelve200ConHXRedirect(t *testing.T) {
	invocado := false
	protegido := ui.Requiere(domain.PermVerPanel, handlerCentinela(&invocado))

	req := httptest.NewRequest(http.MethodGet, "/admin/panel", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	protegido(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "/login?next=%2Fadmin%2Fpanel", rec.Header().Get("HX-Redirect"))
	assert.Empty(t, rec.Header().Get("Location"))
	assert.False(t, invocado)
}

func TestRequiere_rolConPermiso_invocaHandler(t *testing.T) {
	for _, rol := range []domain.Rol{domain.RolAdmin, domain.RolLaboratorio} {
		invocado := false
		protegido := ui.Requiere(domain.PermCargarDiag, handlerCentinela(&invocado))

		req := requestConUsuario(http.MethodGet, "/diagnosticos/alta", &domain.Usuario{ID: 1, Usuario: "u", Rol: rol})
		rec := httptest.NewRecorder()

		protegido(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "rol %q", rol)
		assert.True(t, invocado, "rol %q debe poder cargar diagnosticos", rol)
	}
}

func TestRequiere_autenticadoSinPermiso_devuelve403(t *testing.T) {
	invocado := false
	protegido := ui.Requiere(domain.PermGestionUsuarios, handlerCentinela(&invocado))

	lab := &domain.Usuario{ID: 2, Usuario: "lab", Rol: domain.RolLaboratorio}
	req := requestConUsuario(http.MethodGet, "/admin/usuarios", lab)
	rec := httptest.NewRecorder()

	protegido(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.False(t, invocado)
	assert.Contains(t, rec.Body.String(), "No tiene permisos")
}

func TestConUsuario_cookieValida_dejaUsuarioEnContext(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret-de-prueba")
	auth := application.NewAuthService(new(mocks.MockUsuarioRepository))
	token, err := auth.EmitirToken(&domain.Usuario{ID: 9, Usuario: "admin", Nombre: "Admin", Rol: domain.RolAdmin})
	assert.NoError(t, err)

	var visto *domain.Usuario
	envuelto := ui.ConUsuario(auth, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		visto = domain.UsuarioDeCtx(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: ui.NombreCookieSesion, Value: token})
	rec := httptest.NewRecorder()

	envuelto.ServeHTTP(rec, req)

	assert.NotNil(t, visto)
	assert.True(t, visto.Autenticado())
	assert.Equal(t, domain.RolAdmin, visto.Rol)
	assert.Equal(t, "admin", visto.Usuario)
}

// Una cookie manipulada no debe romper la request: se borra y se sigue como publico.
func TestConUsuario_cookieBasura_borraCookieYSigueComoPublico(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret-de-prueba")
	auth := application.NewAuthService(new(mocks.MockUsuarioRepository))

	var visto *domain.Usuario
	envuelto := ui.ConUsuario(auth, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		visto = domain.UsuarioDeCtx(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: ui.NombreCookieSesion, Value: "token.manipulado.xxx"})
	rec := httptest.NewRecorder()

	envuelto.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "no debe cortar con 500")
	assert.NotNil(t, visto)
	assert.False(t, visto.Autenticado())

	cookie := cookieDeRespuesta(rec, ui.NombreCookieSesion)
	assert.NotNil(t, cookie, "debe borrar la cookie invalida")
	assert.True(t, cookie.MaxAge < 0)
}

func TestConUsuario_sinCookie_sigueComoPublico(t *testing.T) {
	auth := application.NewAuthService(new(mocks.MockUsuarioRepository))

	var visto *domain.Usuario
	envuelto := ui.ConUsuario(auth, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		visto = domain.UsuarioDeCtx(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	envuelto.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotNil(t, visto, "UsuarioDeCtx nunca devuelve nil")
	assert.False(t, visto.Autenticado())
	assert.Nil(t, cookieDeRespuesta(rec, ui.NombreCookieSesion), "sin cookie no hay nada que borrar")
}

func cookieDeRespuesta(rec *httptest.ResponseRecorder, nombre string) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == nombre {
			return c
		}
	}
	return nil
}
