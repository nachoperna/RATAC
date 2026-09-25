package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/views"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

type AdminHandler struct {
	adminService *application.AdminService
	pacienteService *application.PacienteService
	authService *application.AuthService
}

func NewAdminHandler(adminService *application.AdminService, pacienteService *application.PacienteService, authService *application.AuthService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		pacienteService: pacienteService,
		authService: authService,
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

	token, err := r.Cookie(NOMBRE_TOKEN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	email, err := h.authService.GetEmail(r.Context(), token.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, archivo := range archivos {
		contenido, err := archivo.Open() // abrimos el archivo
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer contenido.Close() // cerramos el archivo luego de usarlo
		paciente, err := h.adminService.ConvertirDocumento(contenido, archivo.Filename, email, r.Context())
		if err != nil {
			nombre, _, _ := strings.Cut(filepath.Base(archivo.Filename), ".")
			_ = h.adminService.BorrarTemporal(nombre, nil)
			w.WriteHeader(http.StatusOK)
			views.ErrorCargaDiagnostico("", "El diagnóstico subido ya se encuentra cargado en el sistema.").Render(r.Context(), w)
			return
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
		err = h.adminService.GuardarImagenes(imagenes)
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

func (h *AdminHandler) ShowAdminPanel(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie(NOMBRE_TOKEN)
	if err != nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pacientes, total, err := h.adminService.GetUltimosDiagnosticosCargados(r.Context(), 0, token.Value)
	if err != nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// SE DEBE HACER PETICION A LA BD PARA OBTENER ROL REAL CON EL TOKEN
	solicitudes, autorizado, err := h.adminService.GetSolicitudes(r.Context(), "admin")
	if err != nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var lista_solicitudes templ.Component
	if autorizado {
		lista_solicitudes = views.Solicitudes(solicitudes)
	}
	var tabla_diagnosticos templ.Component = views.TablaUltimosDiagnosticos(pacientes, 0, total, true)
	var paginacion templ.Component = views.ResultadosRestantes(int8(len(pacientes)), 0, total)
	var header templ.Component = views.HeaderLinks(true, "")

	tmp_aux := template.New("panel_administrador.html").Funcs(template.FuncMap{ "render": renderTempl })
	tmp, err := tmp_aux.ParseFiles("./infrastructure/UI/static/panel_administrador.html")
	if err != nil {
		fmt.Printf("Error al parsear el template: %v", err)
		http.Error(w, "No se pudo cargar la página", http.StatusInternalServerError)
		return
	}

	datos := map[string]any{
		"TablaDiagnosticos": tabla_diagnosticos,
		"Paginacion": paginacion,
		"Header": header,
		"Solicitudes": lista_solicitudes,
	}
	tmp.Execute(w, datos)
}

func (h *AdminHandler) DiagnosticosByUser(w http.ResponseWriter, r *http.Request)  {
	offset, err := getOffset(r.URL.Query().Get("offset"))
	if err != nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, err := r.Cookie(NOMBRE_TOKEN)
	if err != nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pacientes, total, err := h.adminService.GetUltimosDiagnosticosCargados(r.Context(), offset, token.Value)
	if err != nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	views.TablaUltimosDiagnosticos(pacientes, offset, total, false).Render(r.Context(), w)
}

func (h *AdminHandler) SolicitudAprobada (w http.ResponseWriter, r *http.Request)  {
	email := r.URL.Query().Get("email")
	nombre, _, err := h.adminService.SolicitudAprobada(r.Context(), email)
	if err != nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	contraseña := nombre + os.Getenv("CONTRA_AUTOGENERADA")
	err = enviarMailSolicitudAprobada(nombre, email, contraseña)
	if err != nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) SolicitudRechazada (w http.ResponseWriter, r *http.Request)  {
	email := r.URL.Query().Get("email")
	nombre, _, err := h.adminService.SolicitudRechazada(r.Context(), email)
	if err != nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = enviarMailSolicitudRechazada(nombre, email)
	if err != nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func mapearCamposAPaciente(info InformacionDiagnostico) domain.Paciente {
	f := info.Categorias
	ref, _ := strconv.ParseBool(f["f-mastocitomas"])
	return domain.Paciente{
		Protocolo:               f["f-protocolo"],
		Fecha:                   f["f-fecha"],
		Solicitante:             f["f-solicitante"],
		Tecnica:                 f["f-tecnica"],
		Familia:                 new(f["f-familia"]),
		Especie:                 new(f["f-especie"]),
		Raza:                    new(f["f-raza"]),
		Edad:                    new(f["f-edad"]),
		NombrePaciente:          f["f-paciente"],
		ReferenciasMastocitomas: ref,
		Antecedentes:            new(f["f-antecedentes"]),
		DescripcionMacroscopica: new(f["f-macroscopica"]),
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
