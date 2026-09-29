package application

import (
	"RATAC/domain"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var tipos_validos = map[string]bool{
	"application/pdf": true,
	"application/zip": true,
}

var (
	ErrTipoArchivoInvalido = errors.New("el archivo no es un PDF ni un DOCX")
	ErrDocumentoIlegible   = errors.New("no se pudo leer el contenido del archivo")
	ErrFormatoDesconocido  = errors.New("el documento no tiene el formato de diagnostico esperado")
	ErrDiagnosticoRepetido = errors.New("el diagnostico ya esta registrado en el sistema")
)

type AdminService struct {
	adminRepo domain.AdminRepository
}

func NewAdminService(adminRepo domain.AdminRepository) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

func (s *AdminService) ConvertirDocumento(archivo multipart.File, nombre, email string, ctx context.Context) (*domain.Paciente, error) {
	buffer := make([]byte, 512) // necesitamos generar un pequeño buffer en memoria RAM de 512 BYTES para leer los primeros bytes del archivo
	_, err := archivo.Read(buffer)
	if err != nil {
		return nil, ErrDocumentoIlegible
	}
	archivo.Seek(0, 0)
	tipo := http.DetectContentType(buffer) // detectamos el tipo del archivo segun su contenido
	if !tipos_validos[tipo] {
		return nil, ErrTipoArchivoInvalido
	}

	os.MkdirAll("./ArchivosTemporales", os.ModePerm)
	tmpFile, err := os.CreateTemp("./ArchivosTemporales/", fmt.Sprintf("TEMP_%s", nombre))
	if err != nil {
		return nil, errors.New("Error al crear archivo temporal")
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, archivo)
	if err != nil {
		return nil, errors.New("Error al copiar contenido a archivo temporal")
	}
	tmpFile.Close()

	// Aca se debe llamar a ejecucion de diag_to_json.py / pdf_to_json.py
	// cmd := exec.Command("docker", "compose", "exec", "app", "python3", "ProcesadoJsons/diag_to_json.py", tmpFile.Name())
	cmd := exec.Command("python3", "ProcesadoJsons/diag_to_json.py", tmpFile.Name(), nombre, email)
	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil { // significa que python termino su ejecucion con un codigo de salida distinto a 0
		return nil, fmt.Errorf("%w: %s", ErrFormatoDesconocido, strings.TrimSpace(stderr.String()))
	}

	// Pasamos el json obtenido a la capa de infraestructura para que lo mapee a nuestro objeto y lo retorne
	paciente, err := s.adminRepo.MapeoDocumento(stdout)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFormatoDesconocido, err)
	}
	// Obtenemos el protocolo del paciente y buscamos si ya estaba registrado en nuestro sistema
	if s.adminRepo.PacienteYaRegistrado(ctx, paciente.Protocolo) {
		return nil, ErrDiagnosticoRepetido
	}
	return paciente, nil
}

func (s *AdminService) BorrarTemporal(archivo string, imagenes []string) error {
	err := os.Remove(fmt.Sprintf("JSONS/TEMP_%s.json", archivo))
	if err != nil {
		return err
	}

	for _, img := range imagenes {
		err := os.Remove(strings.TrimPrefix(img, "/"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *AdminService) RenombrarTemporal(nombre string) error {
	json, err := filepath.Glob(fmt.Sprintf("JSONS/*%s*.json", nombre))
	if err != nil {
		return err
	}
	if json == nil {
		return errors.New("No se encontro el archivo json")
	}

	err = os.Rename(json[0], fmt.Sprintf("JSONS/%s.json", nombre))
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminService) GetImagenesHuerfanas(imagenes []string, nombre string) ([]string, error) {
	imgs, err := filepath.Glob(fmt.Sprintf("IMAGENES/*%s*", nombre))
	if err != nil {
		return nil, err
	}

	var huerfanas []string
	for _, img := range imgs {
		// La imagen es huerfana si no pertenece a la informacion de diagnostico final
		if !slices.Contains(imagenes, "/"+img) {
			huerfanas = append(huerfanas, img)
		}
	}
	return huerfanas, nil
}

func (s *AdminService) GenerarJson(nombre string, paciente domain.Paciente) error {
	archivo, err := os.Create(fmt.Sprintf("JSONS/%s.json", nombre))
	if err != nil {
		return err
	}
	err = json.NewEncoder(archivo).Encode(paciente)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminService) GuardarImagenes(imagenes []*multipart.FileHeader) error {
	for _, img := range imagenes {
		nueva_imagen, err := os.Create(fmt.Sprintf("IMAGENES/%s", img.Filename))
		if err != nil {
			return err
		}
		contenido, err := img.Open()
		if err != nil {
			return err
		}
		defer contenido.Close()
		_, err = io.Copy(nueva_imagen, contenido)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *AdminService) GetUltimosDiagnosticosCargados(ctx context.Context, offset int8, token string) ([]domain.Paciente, int16, error) {
	diagnosticos, total, err := s.adminRepo.GetUltimosDiagnosticosCargados(ctx, offset, token)
	if err != nil {
		return nil, 0, err
	}
	return diagnosticos, total, nil
}

func (s *AdminService) GetSolicitudes(ctx context.Context, rol string) ([]domain.Solicitud, bool, error) {
	if rol != "admin"{
		return nil, false, nil
	}
	solicitudes, err := s.adminRepo.GetSolicitudes(ctx)
	if err != nil {
		return nil, true, err
	}
	return solicitudes, true, nil
}

func (s *AdminService) SolicitudAprobada (ctx context.Context, email string) (string, string, error) {
	return s.adminRepo.SolicitudAprobada(ctx, email)
}

func (s *AdminService) SolicitudRechazada (ctx context.Context, email string) (string, error) {
	return s.adminRepo.SolicitudRechazada(ctx, email)
}
