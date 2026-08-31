package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/views"
	"net/http"
	"os"
	"strings"
)

const destinoPorDefecto = "/admin/panel"

type AuthHandler struct {
	authService *application.AuthService
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	if domain.UsuarioDeCtx(r.Context()).Autenticado() {
		http.Redirect(w, r, destinoPorDefecto, http.StatusSeeOther)
		return
	}
	next := destinoSeguro(r.URL.Query().Get("next"))
	render(w, r, views.Login(next, ""), false)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}
	usuario := strings.TrimSpace(r.FormValue("usuario"))
	password := r.FormValue("password")
	next := destinoSeguro(r.FormValue("next"))

	u, err := h.authService.Login(r.Context(), usuario, password)
	if err != nil {
		// 200 y no 401: HTMX por defecto no swapea el body de respuestas 4xx,
		// y el mensaje de error no se mostraria.
		w.WriteHeader(http.StatusOK)
		views.LoginFormulario(next, "Usuario o contraseña incorrectos").Render(r.Context(), w)
		return
	}

	token, err := h.authService.EmitirToken(u)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		views.LoginFormulario(next, "No se pudo iniciar sesión, intente de nuevo").Render(r.Context(), w)
		return
	}
	fijarCookieSesion(w, token)

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", next)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	borrarCookieSesion(w)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func fijarCookieSesion(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     NombreCookieSesion,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure fijo en true impediria el login sobre http://localhost:8080.
		Secure: os.Getenv("APP_ENV") == "production",
		MaxAge: int(application.DuracionSesion.Seconds()),
	})
}

func borrarCookieSesion(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     NombreCookieSesion,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   os.Getenv("APP_ENV") == "production",
		MaxAge:   -1,
	})
}

// destinoSeguro evita un open redirect: solo acepta rutas internas. Rechaza las
// que empiezan con "//" (que el navegador interpreta como host externo) y las
// que no arrancan con "/".
func destinoSeguro(next string) string {
	if next == "" || !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		return destinoPorDefecto
	}
	if strings.Contains(next, "\\") || strings.ContainsAny(next, "\r\n") {
		return destinoPorDefecto
	}
	return next
}
