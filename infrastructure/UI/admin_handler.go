package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/views"
	"encoding/json"
	"fmt"
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

type DescMicro struct {
    Descripcion string   `json:"descripcion"`
    Diagnostico string   `json:"diagnostico"`
    Imagenes    []string `json:"imagenes"`
}

type InformacionDiagnostico struct {
	Categorias	map[string]string	`json:"fields"`
	DescMicros	[]DescMicro		`json:"microCards"`
	Imagenes	[]string		`json:"images"`
}

type PayloadDiagnostico struct {
	Archivo	string			`json:"archivo"`
	Cambios	string			`json:"cambios"`
	Campos	InformacionDiagnostico	`json:"campos"`
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
	r.ParseMultipartForm(32 << 20)
	
	nombre_original := r.FormValue("archivo")
	nombre_base, _, _ := strings.Cut(filepath.Base(nombre_original), ".")
	cambios, err := strconv.ParseBool(r.FormValue("cambios"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !cambios{ // significa que ya se encuentra en el servidor el json temporal y las imagenes temporales 
		// renombrar json quitando el "TEMP_" del nombre de archivo 
		err = h.adminService.RenombrarTemporal(nombre_base)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.pacienteService.InsertarDiagnostico(r.Context(), nombre_base)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Se necesitan renders de errores correspondientes
	} else { // Debemos reprocesar todos los valores de los campos (mas facil que hacer una especie de diff y editar ael json)
		// BORRAMOS json temporal e imagenes no incluidas en el diagnostico final
		var info_diagnostico InformacionDiagnostico
		err = json.Unmarshal([]byte(r.FormValue("campos")), &info_diagnostico)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		huerfanas, err := h.adminService.GetImagenesHuerfanas(info_diagnostico.Imagenes, nombre_original)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = h.adminService.BorrarTemporal(nombre_base, huerfanas)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		paciente := mapearCamposAPaciente(info_diagnostico)
		err = h.adminService.GenerarJson(nombre_base, paciente)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imagenes := r.MultipartForm.File["imagenes"]
		err = h.adminService.GuardarImagenes(imagenes, nombre_original)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = h.pacienteService.InsertarDiagnostico(r.Context(), nombre_base)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	// Render con opcion de redireccion a la pagina completa del diagnostico
}

func mapearCamposAPaciente(info InformacionDiagnostico) domain.Paciente {
	f := info.Categorias
	ref, _ := strconv.ParseBool(f["f-mastocitomas"])
	return domain.Paciente{
		Protocolo:               f["f-protocolo"],
		Fecha:                   f["f-fecha"],
		Solicitante:             f["f-solicitante"],
		Tecnica:                 f["f-tecnica"],
		Familia:                 strPointer(f["f-propietario"]),
		Especie:                 strPointer(f["f-especie"]),
		Raza:                    strPointer(f["f-raza"]),
		Edad:                    strPointer(f["f-edad"]),
		NombrePaciente:          f["f-paciente"],
		ReferenciasMastocitomas: ref,
		Antecedentes:            strPointer(f["f-antecedentes"]),
		DescripcionMacroscopica: strPointer(f["f-macroscopica"]),
		Descripciones_microscopicas: mapearDescMicros(info.DescMicros),
	}
}

func mapearDescMicros(cards []DescMicro) []domain.Descripcion_microscopicas {
    var descripciones []domain.Descripcion_microscopicas
    for _, c := range cards {
        d := c.Diagnostico
        descripciones = append(descripciones, domain.Descripcion_microscopicas{
            Descripcion: c.Descripcion,
            Diagnostico: domain.Diagnostico{
                Descripcion: 	&d,
                Imagenes:	rutaCorrectaImagenes(c.Imagenes),
            },
            TablaGrado: []domain.Grado_oncologico{},
        })
    }
    return descripciones
}

func rutaCorrectaImagenes(imagenes []string) []string {
	var aux []string
	for _, img := range imagenes {
		if strings.HasPrefix(img, "/IMAGENES/") {
			aux = append(aux, strings.TrimPrefix(img, "/"))
		} else {
			aux = append(aux, fmt.Sprintf("IMAGENES/%s", img))
		}
	}
	return aux
}
func strPointer(s string) *string { return &s }
