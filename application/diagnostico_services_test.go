package application_test

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDiagnostico_parametrosValidos_retornaDiagnostico(t *testing.T) {
	mockRepo := new(mocks.MockDiagnosticoRepository)
	service := application.NewDiagnosticoService(mockRepo)

	protocolo := "1234-A"
	descMicro := "Descripcion de prueba"
	descTexto := "Diagnostico positivo"
	esperado := &domain.Diagnostico{
		Descripcion: &descTexto,
		Imagenes:    []string{"img1.png", "img2.png"},
	}

	mockRepo.On("GetDiagnostico", mock.Anything, protocolo, descMicro).Return(esperado, nil)

	resultado, err := service.GetDiagnostico(context.Background(), protocolo, descMicro)

	assert.NoError(t, err)
	assert.Equal(t, esperado, resultado)
	mockRepo.AssertExpectations(t)
}

func TestGetDiagnostico_errorEnBaseDeDatos_retornaError(t *testing.T) {
	mockRepo := new(mocks.MockDiagnosticoRepository)
	service := application.NewDiagnosticoService(mockRepo)

	errorDB := errors.New("conexion fallida")
	mockRepo.On("GetDiagnostico", mock.Anything, "1234-A", "desc").Return(nil, errorDB)

	resultado, err := service.GetDiagnostico(context.Background(), "1234-A", "desc")

	assert.Error(t, err)
	assert.Equal(t, errorDB, err)
	assert.Nil(t, resultado)
	mockRepo.AssertExpectations(t)
}
func TestCreateDiagnostico_parametrosValidos_retornaNil(t *testing.T) {
	mockRepo := new(mocks.MockDiagnosticoRepository)
	service := application.NewDiagnosticoService(mockRepo)

	protocolo := "CAN-001"
	descMicro := "Descripcion de prueba"
	nuevoDiag := domain.Diagnostico{
		Imagenes: []string{"nueva_img.png"},
	}

	mockRepo.On("CreateDiagnostico", mock.Anything, protocolo, descMicro, nuevoDiag).Return(nil)

	err := service.CreateDiagnostico(context.Background(), protocolo, descMicro, nuevoDiag)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteImagen_parametrosValidos_retornaNil(t *testing.T) {
	mockRepo := new(mocks.MockDiagnosticoRepository)
	service := application.NewDiagnosticoService(mockRepo)

	rutaImagen := "IMAGENES/img_borrar.png"

	mockRepo.On("DeleteImagen", mock.Anything, rutaImagen).Return(nil)

	err := service.DeleteImagen(context.Background(), rutaImagen)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCountImagenes_ejecucionExitosa_retornaCantidad(t *testing.T) {
	mockRepo := new(mocks.MockDiagnosticoRepository)
	service := application.NewDiagnosticoService(mockRepo)

	var cantidadEsperada int64 = 15

	mockRepo.On("CountImagenes", mock.Anything).Return(cantidadEsperada, nil)

	resultado, err := service.CountImagenes(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, cantidadEsperada, resultado)
	mockRepo.AssertExpectations(t)
}
