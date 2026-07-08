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
	tipo := http.DetectContentType(buffer) // detectamos el tipo del archivo segun su contenido
	if !tipos_validos[tipo] {
		return nil, err
	}
	
	tmpFile, _ := os.CreateTemp("./ArchivosTemporales/", fmt.Sprintf("TEMP_%s", nombre))
	defer os.Remove(tmpFile.Name())
	io.Copy(tmpFile, archivo)
	tmpFile.Close()

	// Aca se debe llamar a ejecucion de diag_to_json.py / pdf_to_json.py
	cmd := exec.Command("python3", "../ProcesadoJsons/diag_to_json.py", tmpFile.Name())
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
