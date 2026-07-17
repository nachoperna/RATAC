package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/views"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type AdminHandler struct {
	adminService *application.AdminService
	pacienteService *application.PacienteService
}

func NewAdminHandler(adminService *application.AdminService, pacienteService *application.PacienteService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		pacienteService: pacienteService,
	}
}

func (h *AdminHandler) ProcesarDocumento(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	archivos := r.MultipartForm.File["archivos"]
	var pacientes []domain.Paciente

	for _, archivo := range archivos {
		contenido, err := archivo.Open() // abrimos el archivo
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer contenido.Close() // cerramos el archivo luego de usarlo
		paciente, err := h.adminService.ConvertirDocumento(contenido, archivo.Filename)
		if err != nil {
			// w.WriteHeader(http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
			// Renderizar templ de error
		} else {
			pacientes = append(pacientes, *paciente)
		}
	}

	w.WriteHeader(http.StatusOK)
	views.InformacionExtraida(pacientes).Render(r.Context(), w)
}

func (h *AdminHandler) BorrarTemporal(w http.ResponseWriter, r *http.Request)  {
	archivo := r.URL.Query().Get("archivo")
	if archivo == ""{
		w.WriteHeader(200)
		return
	}
	var imagenes []string = r.URL.Query()["imagenes"]

	nombre, _, _ := strings.Cut(filepath.Base(archivo), ".")
	err := h.adminService.BorrarTemporal(nombre, imagenes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) AltaDiagnostico(w http.ResponseWriter, r *http.Request)  {
	nombre, _, _ := strings.Cut(filepath.Base(r.FormValue("archivo")), ".")
	cambios, err := strconv.ParseBool(r.FormValue("cambios"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !cambios{ // significa que ya se encuentra en el servidor el json temporal y las imagenes temporales 
		// renombrar json quitando el "TEMP_" del nombre de archivo 
		err = h.adminService.RenombrarTemporal(nombre)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.pacienteService.InsertarDiagnostico(r.Context(), nombre)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Se necesitan renders de errores correspondientes
	}
	w.WriteHeader(http.StatusOK)
	// Render con opcion de redireccion a la pagina completa del diagnostico
}
