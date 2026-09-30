package ui

import (
	"RATAC/application"
	"RATAC/views"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AuthHandler struct {
	AuthService  *application.AuthService
	AdminService *application.AdminService
}

func NewAuthHandler(AuthService *application.AuthService, AdminService *application.AdminService) *AuthHandler {
	return &AuthHandler{
		AuthService:  AuthService,
		AdminService: AdminService,
	}
}

const NOMBRE_TOKEN = "token_sesion"

func (h *AuthHandler) LoggerChecker(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(NOMBRE_TOKEN)
		if err != nil || token.Value == "" {
			h.ifCargaFallida(r)
			setRedireccionIngreso(w)
			return
		}
		if !r.URL.Query().Has("login_reciente") { // Evito la consulta si el usuario accedio a la ruta desde una redireccion por login/register
			activa, err := h.AuthService.Validacion(r.Context(), token.Value)
			if err != nil || !activa {
				h.ifCargaFallida(r)
				setRedireccionIngreso(w)
				return
			}
		}
		handler(w, r)
	}
}

func setRedireccionIngreso(w http.ResponseWriter) {
	w.Header().Set("HX-Redirect", "/ingreso/formulario")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	contenido, _ := os.ReadFile("./infrastructure/UI/static/acceso_restringido.html")
	w.Write(contenido)
}

func (h *AuthHandler) ifCargaFallida(r *http.Request) {
	if r.URL.Path == "/diagnosticos/alta/carga" { // borramos los archivos temporales creados
		nombre_base, _, _ := strings.Cut(filepath.Base(r.FormValue("archivo")), ".")
		imgs, _ := h.AdminService.GetImagenesHuerfanas([]string{}, nombre_base)
		go h.AdminService.BorrarTemporal(nombre_base, imgs) // borramos en segundo plano
		// en el futuro se debe tratar el error de forma asincrona para terminar de borrar todos los archivos temporales
	}
}

func (h *AuthHandler) NoActionOk(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	w.WriteHeader(http.StatusOK)
}
func (h *AuthHandler) NoActionErr(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	w.WriteHeader(http.StatusBadRequest)
}

// responderRegistroFallido devuelve el motivo del rechazo como texto plano, que
// el modal de solicitud muestra en su cuerpo (ver selectModal en panel_acceso.js).
// Se mantiene el estado de error para que htmx lo trate como tal.
func responderRegistroFallido(w http.ResponseWriter, mensaje string, detalle error, contexto string, inicio time.Time) {
	if detalle != nil {
		log.Printf("[%s] %v", contexto, detalle)
	}
	// Un correo ya registrado consulta la base y tarda mas que un rechazo por
	// campos vacios; sin este piso la diferencia delata que correos tienen cuenta.
	restante := DURACION_MINIMA_LOGIN - time.Since(inicio)
	if restante > 0 {
		time.Sleep(restante)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(mensaje))
}

func (h *AuthHandler) Registrarse(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	err := r.ParseForm()
	if err != nil {
		responderRegistroFallido(w, MSJ_ERROR_REGISTRO, err, "Registrarse: parsear formulario", inicio)
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	nombre_lab := strings.TrimSpace(r.FormValue("nombre-lab"))
	contraseña := nombre_lab + os.Getenv("CONTRA_AUTOGENERADA") // Luego el laboratorio debera cambiarla
	veterinarios := r.Form["nombre-vet"]
	matriculas := r.Form["matricula-vet"]
	ciudad_origen := strings.TrimSpace(r.FormValue("ciudad-origen"))

	err = h.AuthService.Registrarse(r.Context(), email, contraseña, ciudad_origen, nombre_lab, matriculas, veterinarios)
	if err != nil {
		// Se responde lo mismo haya fallado por un correo repetido o por cualquier
		// otra causa: distinguirlos permitiria averiguar que correos tienen cuenta.
		responderRegistroFallido(w, MSJ_ERROR_REGISTRO_EMAIL, err, "Registrarse: "+email, inicio)
		return
	}
	go enviarMailSolicitud(nombre_lab, email, ciudad_origen, veterinarios, matriculas)
	w.WriteHeader(http.StatusOK)
}

// DURACION_MINIMA_LOGIN es el piso de tiempo que tarda todo intento de acceso
// fallido. Sin el, un email inexistente responde en ~2ms (falla la consulta y
// nunca se ejecuta bcrypt) mientras que un email real con la contraseña
// equivocada tarda ~100ms, lo que permitiria descubrir que usuarios estan
// registrados midiendo la respuesta.
const DURACION_MINIMA_LOGIN = 250 * time.Millisecond

// responderLoginFallido devuelve el aviso generico de credenciales invalidas.
// Se usa la misma respuesta para cualquier causa de fallo (usuario inexistente,
// contraseña incorrecta o error al cambiarla) y se espera hasta completar
// DURACION_MINIMA_LOGIN para que todos los intentos tarden lo mismo.
func responderLoginFallido(w http.ResponseWriter, r *http.Request, inicio time.Time) {
	if restante := DURACION_MINIMA_LOGIN - time.Since(inicio); restante > 0 {
		time.Sleep(restante)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Se responde 200 porque htmx descarta el cuerpo de las respuestas de error y
	// el cartel nunca llegaria a mostrarse en el formulario.
	w.WriteHeader(http.StatusOK)
	views.ErrorLogin().Render(r.Context(), w)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	email := r.FormValue("email")
	contraseña := r.FormValue("contraseña")
	nueva_contraseña := r.FormValue("nueva_contraseña")
	if nueva_contraseña != "" {
		err := h.AuthService.CambiarContraseña(r.Context(), email, contraseña, nueva_contraseña)
		if err != nil {
			responderLoginFallido(w, r, inicio)
			return
		}
		contraseña = nueva_contraseña
	}
	token, expiracion, err := h.AuthService.Login(r.Context(), email, contraseña)
	if err != nil {
		responderLoginFallido(w, r, inicio)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     NOMBRE_TOKEN,
		Value:    token,
		Expires:  expiracion,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	w.Header().Set("HX-Redirect", "/admin/panel?login_reciente=true")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie(NOMBRE_TOKEN)
	if err != nil {
		http.Error(w, "Error al desloguear", http.StatusBadRequest)
		return
	}

	_ = h.AuthService.Logout(r.Context(), token.Value) // si hay error en borrar la sesion de la bd deslogueamos igual

	http.SetCookie(w, &http.Cookie{
		Name:     NOMBRE_TOKEN,
		Value:    "",
		Expires:  time.Now(),
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) SesionActiva(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie(NOMBRE_TOKEN)
	if err != nil {
		http.Error(w, "SESION TERMINADA", http.StatusBadRequest)
		return
	}
	activa, err := h.AuthService.Validacion(r.Context(), token.Value)
	if err != nil || !activa {
		http.Error(w, "SESION TERMINADA", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1> SESION ACTIVA </h1>"))
}
