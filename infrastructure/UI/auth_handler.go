package ui

import (
	"RATAC/application"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AuthHandler struct {
	AuthService *application.AuthService
	AdminService *application.AdminService
}

func NewAuthHandler(AuthService *application.AuthService, AdminService *application.AdminService) *AuthHandler {
	return &AuthHandler{
		AuthService: AuthService,
		AdminService: AdminService,
	}
}

const NOMBRE_TOKEN = "token_sesion"

func (h *AuthHandler) LoggerChecker(handler http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(NOMBRE_TOKEN)
		if err != nil || token.Value == ""{
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
		handler(w,r)
	}
}

func setRedireccionIngreso(w http.ResponseWriter)  {
	w.Header().Set("HX-Redirect", "/ingreso/formulario")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	contenido, _ := os.ReadFile("./infrastructure/UI/static/acceso_restringido.html")
	w.Write(contenido)
}

func (h *AuthHandler) ifCargaFallida (r *http.Request)  {
	if r.URL.Path == "/diagnosticos/alta/carga" { // borramos los archivos temporales creados
		nombre_base, _, _ := strings.Cut(filepath.Base(r.FormValue("archivo")), ".")
		imgs, _ := h.AdminService.GetImagenesHuerfanas([]string{}, nombre_base)
		go h.AdminService.BorrarTemporal(nombre_base, imgs) // borramos en segundo plano
		// en el futuro se debe tratar el error de forma asincrona para terminar de borrar todos los archivos temporales
	}
}

func (h *AuthHandler) Registrarse (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// renderizar templ de error
		http.Error(w, "Error al parsear formulario: " + err.Error(), http.StatusBadRequest)
		return 
	}
	email := r.FormValue("email")
	nombre_lab := r.FormValue("nombre-lab")
	contraseña := fmt.Sprintf("%s-RATAC-2026", nombre_lab) // Luego el laboratorio debera cambiarla
	veterinarios := r.Form["nombre-vet"]
	matriculas := r.Form["matricula-vet"]
	ciudad_origen := r.FormValue("ciudad-origen")

	err = h.AuthService.Registrarse(r.Context(), email, contraseña, ciudad_origen, nombre_lab, matriculas, veterinarios)
	if err != nil {
		// renderizar templ de error
		http.Error(w, "Error al iniciar seison" + err.Error(), http.StatusBadRequest)
		return 
	}
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
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path: "/",
	})
	w.Header().Set("HX-Redirect", "/admin/panel?login_reciente=true")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout (w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie(NOMBRE_TOKEN)
	if err != nil {
		http.Error(w, "Error al desloguear", http.StatusBadRequest)
		return
	}

	_ = h.AuthService.Logout(r.Context(), token.Value) // si hay error en borrar la sesion de la bd deslogueamos igual
	
	http.SetCookie(w, &http.Cookie{
		Name:    NOMBRE_TOKEN,
		Value:   "",
		Expires: time.Now(),
		MaxAge:  -1,
		Path:    "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
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
