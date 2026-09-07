package ui

import (
	"RATAC/application"
	"net/http"
	"time"
)

type AuthHandler struct {
	AuthService *application.AuthService
}

func NewAuthHandler(AuthService *application.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: AuthService}
}

const NOMBRE_TOKEN = "token_sesion"

func (h *AuthHandler) Registrarse (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// renderizar templ de error
		http.Error(w, "Error al parsear formulario: " + err.Error(), http.StatusBadRequest)
		return 
	}
	email := r.FormValue("email")
	nombre_lab := r.FormValue("nombre-lab")
	contraseña := "contraseña-auto-generada" // Luego el laboratorio debera cambiarla
	// veterinarios := r.Form["nombre-vet"]
	// matriculas := r.Form["matricula-vet"]
	ciudad_origen := r.FormValue("ciudad-origen")
	rol := "laboratorio"

	token, expiracion, err := h.AuthService.Registrarse(r.Context(), email, contraseña, rol, ciudad_origen, nombre_lab)
	if err != nil {
		// renderizar templ de error
		http.Error(w, "Error al iniciar seison" + err.Error(), http.StatusBadRequest)
		return 
	}
	http.SetCookie(w, &http.Cookie{
		Name: NOMBRE_TOKEN,
		Value: token,
		Expires: expiracion,
	})
	w.Header().Set("HX-Redirect", "/admin/panel")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Login (w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	contraseña := r.FormValue("contraseña")
	nueva_contraseña := r.FormValue("nueva_contraseña")
	if (nueva_contraseña != ""){
		err := h.AuthService.CambiarContraseña(r.Context(), email, contraseña, nueva_contraseña)
		if err != nil {
			// renderizar templ de error
			http.Error(w, "Error al cambiar contraseña", http.StatusBadRequest)
			return 
		}
		contraseña = nueva_contraseña
	}
	token, expiracion, err := h.AuthService.Login(r.Context(), email, contraseña)
	if err != nil {
		// renderizar templ de error
		http.Error(w, "Error al iniciar seison", http.StatusBadRequest)
		return 
	}
	http.SetCookie(w, &http.Cookie{
		Name: NOMBRE_TOKEN,
		Value: token,
		Expires: expiracion,
	})
	w.Header().Set("HX-Redirect", "/admin/panel")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout (w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie(NOMBRE_TOKEN)
	if err != nil {
		http.Error(w, "SESION TERMINADA", http.StatusBadRequest)
		return
	}

	err = h.AuthService.Logout(r.Context(), token.Value)
	if err != nil {
		http.Error(w, "Error al desloguear", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    NOMBRE_TOKEN,
		Value:   "",
		Expires: time.Now(),
		MaxAge:  -1,
		Path:    "/",
	})
	w.Write([]byte("<h1> DESLOGUEADO CON EXITO </h1>")) 
}

func (h *AuthHandler) SesionActiva (w http.ResponseWriter, r *http.Request) {
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
