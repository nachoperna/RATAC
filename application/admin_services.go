package application

import (
	"RATAC/domain"
	"bytes"
	"mime/multipart"
	"net/http"
	"os/exec"
)

// type AdminService struct {
// 	adminRepo domain.AdminRepository
// }
//
// func NewAdminService(adminRepo domain.AdminRepository) *AdminService {
// 	return &AdminService{
// 		adminRepo: adminRepo,
// 	}
// }

var tipos_validos = map[string]bool {
	"application/pdf": true,
	"application/zip": true,
}

type AdminService struct{
	adminRepo domain.AdminRepository
}

func (s *AdminService) ConvertirDocumento(archivo multipart.File) (*domain.Paciente, error) {
	buffer := make([]byte, 512) // necesitamos generar un pequeño buffer en memoria RAM de 512 BYTES para leer los primeros bytes del archivo
	_, err := archivo.Read(buffer)
	if err != nil {
		return nil, err
	}
	tipo := http.DetectContentType(buffer) // detectamos el tipo del archivo segun su contenido
	if !tipos_validos[tipo] {
		return nil, err
	}

	// Aca se debe llamar a ejecucion de diag_to_json.py / pdf_to_json.py
	cmd := exec.Command("python3", "../ProcesadoJsons/diag_to_json.py")
	cmd.Stdin = archivo
	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// Pasamos el json obtenido a la capa de infraestructura para que lo mapee a nuestro objeto y lo retorne
	paciente, err := s.adminRepo.MapeoDocumento(stdout)
	if err != nil {
		return nil, err
	}
	return paciente, nil
}
