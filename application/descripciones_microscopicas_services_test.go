package application_test

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDescripcionesByProtocolo_parametrosValidos_retornaDescripciones(t *testing.T) {
	mockRepo := new(mocks.MockDescripcionMicroscopicaRepository)
	service := application.NewDescripcionMicroscopicaService(mockRepo)

	protocolo := "CAN-2026-001"
	esperado := []domain.Descripcion_microscopicas{
		{Descripcion: "Carcinoma de células escamosas"},
		{Descripcion: "Bordes limpios"},
	}

	mockRepo.On("GetDescripcionMicroscopicaByProtocolo", mock.Anything, protocolo).Return(esperado, nil)

	resultado, err := service.GetDescripcionesByProtocolo(context.Background(), protocolo)

	assert.NoError(t, err)
	assert.Equal(t, esperado, resultado)
	assert.Len(t, resultado, 2) // Verificamos que traiga exactamente 2 descripciones
	mockRepo.AssertExpectations(t)
}

func TestCreateDescripcionMicroscopica_parametrosValidos_retornaNil(t *testing.T) {
	mockRepo := new(mocks.MockDescripcionMicroscopicaRepository)
	service := application.NewDescripcionMicroscopicaService(mockRepo)

	protocolo := "CAN-2026-002"
	nuevaDesc := domain.Descripcion_microscopicas{
		Descripcion: "Tejido necrótico severo",
	}

	// Al ser un Create que devuelve un error (o nil si hay éxito), mockeamos que retorne nil
	mockRepo.On("CreateDescripcionMicroscopica", mock.Anything, protocolo, nuevaDesc).Return(nil)

	err := service.CreateDescripcionMicroscopica(context.Background(), protocolo, nuevaDesc)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteDescripcionMicroscopica_parametrosValidos_retornaNil(t *testing.T) {
	mockRepo := new(mocks.MockDescripcionMicroscopicaRepository)
	service := application.NewDescripcionMicroscopicaService(mockRepo)

	protocolo := "CAN-2026-003"
	descABorrar := "Tejido necrótico severo"

	// Mockeamos que el borrado es exitoso (devuelve error nil)
	mockRepo.On("DeleteDescripcionMicroscopica", mock.Anything, descABorrar, protocolo).Return(nil)

	err := service.DeleteDescripcionMicroscopica(context.Background(), descABorrar, protocolo)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCountDiagnosticos_ejecucionExitosa_retornaCantidad(t *testing.T) {
	mockRepo := new(mocks.MockDescripcionMicroscopicaRepository)
	service := application.NewDescripcionMicroscopicaService(mockRepo)

	var cantidadEsperada int64 = 42

	// Mockeamos el Count para que devuelva 42 y ningún error
	mockRepo.On("CountDiagnosticos", mock.Anything).Return(cantidadEsperada, nil)

	resultado, err := service.CountDiagnosticos(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, cantidadEsperada, resultado)
	mockRepo.AssertExpectations(t)
}
