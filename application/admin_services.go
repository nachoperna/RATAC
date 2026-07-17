package application

import (
	"RATAC/domain"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var tipos_validos = map[string]bool {
	"application/pdf": true,
	"application/zip": true,
}

type AdminService struct{
	adminRepo domain.AdminRepository
}

func NewAdminService(adminRepo domain.AdminRepository) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

/*
* POSIBLE LOGICA: 
* 1. python crea el json real y ademas devuelve el contenido al stdout
* 2. adminservice retorna el objeto al handler
* 3. adminhandler renderiza informacion en pantalla
* 	4. si el usuario no edito ningun campo de informacion extraida, el archivo json creado se guarda tal cual en JSONS/
* 	4. si el usuario edito algun campo de informacion extraida se usa un discernible para que el servidor sepa
* 		y se mapean todos los campos a un objeto para luego crear el json 
* 		porque es mas simple hacerlo de vuelta que editar un archivo en una linea especifica
* */

func (s *AdminService) ConvertirDocumento(archivo multipart.File, nombre string) (*domain.Paciente, error) {
	buffer := make([]byte, 512) // necesitamos generar un pequeño buffer en memoria RAM de 512 BYTES para leer los primeros bytes del archivo
	_, err := archivo.Read(buffer)
	if err != nil {
		return nil, err
	}
	archivo.Seek(0,0)
	tipo := http.DetectContentType(buffer) // detectamos el tipo del archivo segun su contenido
	if !tipos_validos[tipo] {
		return nil, err
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
	cmd := exec.Command("python3", "ProcesadoJsons/diag_to_json.py", tmpFile.Name(), nombre)
	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil { // significa que python termino su ejecucion con un codigo de salida distinto a 0
		return nil, errors.New(stderr.String())
	}

	// Pasamos el json obtenido a la capa de infraestructura para que lo mapee a nuestro objeto y lo retorne
	paciente, err := s.adminRepo.MapeoDocumento(stdout)
	if err != nil {
		return nil, err
	}
	return paciente, nil
}

func (s *AdminService) BorrarTemporal(archivo string, imagenes []string) error {
	err := os.Remove(fmt.Sprintf("JSONS/%s", fmt.Sprintf("TEMP_%s.json", archivo)))
	if err != nil {
		return errors.New("Error borrando json temporal")
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
