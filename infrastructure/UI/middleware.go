package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/views"
	"net/http"
	"net/url"
)

const NombreCookieSesion = "ratac_sesion"

// ConUsuario envuelve todo el mux. Lee la cookie de sesion si existe, la valida
// y deja el usuario en el context. Nunca bloquea: si el token es invalido o
// vencio, borra la cookie y sigue como publico.
func ConUsuario(auth *application.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(NombreCookieSesion)
		if err != nil || cookie.Value == "" {
			next.ServeHTTP(w, r)
			return
		}
		usuario, err := auth.ValidarToken(cookie.Value)
		if err != nil || usuario == nil {
			borrarCookieSesion(w)
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(domain.ConUsuarioCtx(r.Context(), usuario)))
	})
}

// Requiere envuelve un handler concreto. Sin sesion redirige al login; con
// sesion pero sin el permiso responde 403.
func Requiere(permiso domain.Permiso, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usuario := domain.UsuarioDeCtx(r.Context())

		if !usuario.Autenticado() {
			redirigirALogin(w, r)
			return
		}
		if !domain.Puede(usuario, permiso) {
			w.WriteHeader(http.StatusForbidden)
			views.SinPermiso().Render(r.Context(), w)
			return
		}
		handler(w, r)
	}
}

func redirigirALogin(w http.ResponseWriter, r *http.Request) {
	destino := "/login?next=" + url.QueryEscape(r.URL.RequestURI())

	// Ante un request de HTMX hay que responder 200 con HX-Redirect: un 302 lo
	// seguiria XHR e inyectaria el HTML del login dentro del target del swap.
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", destino)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, destino, http.StatusSeeOther)
}
