package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/views"
	"net/http"
)

type AdminHandler struct {
	adminService *application.AdminService
}
//
// func NewAdminHandler(adminService *application.AdminService) *AdminHandler {
// 	return &AdminHandler{adminService: adminService}
// }

func (h *AdminHandler) ProcesarDocumento(w http.ResponseWriter, r *http.Request) {
	archivos := r.MultipartForm.File["archivos"]
	var pacientes []domain.Paciente

	for _, archivo := range archivos {
		contenido, err := archivo.Open() // abrimos el archivo
		if err != nil {
			w.Header().Set("HX-Trigger", "failed_open")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		defer contenido.Close() // cerramos el archivo luego de usarlo
		paciente, err := h.adminService.ConvertirDocumento(contenido)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			// Renderizar templ de error
		}

		pacientes = append(pacientes, *paciente)
	}

	w.WriteHeader(http.StatusOK)
	views.InformacionExtraida(pacientes).Render(r.Context(), w)
}
