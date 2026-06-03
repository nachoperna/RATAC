package ui_test

import (
	"RATAC/application"
	"RATAC/domain"
	ui "RATAC/infrastructure/UI"
	"RATAC/mocks"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAPICreateDescripcion_parametrosYBodyValidos_retornaStatusCreated(t *testing.T) {
	mockRepo := new(mocks.MockDescripcionMicroscopicaRepository)
	service := application.NewDescripcionMicroscopicaService(mockRepo)
	handler := ui.NewDescripcionMicroscopicaHandler(service)

	// Le decimos al mock que no devuelva errores al guardar
	mockRepo.On("CreateDescripcionMicroscopica", mock.Anything, "CAN-001", mock.AnythingOfType("domain.Descripcion_microscopicas")).Return(nil)

	// Armamos un JSON válido para simular el body que enviaría el frontend o un script
	jsonBody := []byte(`{"Descripcion": "Mastocitoma grado II"}`)
	bodyReader := bytes.NewBuffer(jsonBody)

	// Creamos el request POST pasando el bodyReader
	req := httptest.NewRequest(http.MethodPost, "/api/descripciones?protocolo=CAN-001", bodyReader)
	req.Header.Set("Content-Type", "application/json") // Buena práctica agregar el header

	rec := httptest.NewRecorder()
	handler.APICreateDescripcion(rec, req)

	// Afirmaciones
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAPICreateDescripcion_jsonInvalido_retornaBadRequest(t *testing.T) {
	mockRepo := new(mocks.MockDescripcionMicroscopicaRepository)
	service := application.NewDescripcionMicroscopicaService(mockRepo)
	handler := ui.NewDescripcionMicroscopicaHandler(service)

	// Armamos un JSON malformado (falta cerrar la llave)
	jsonBody := []byte(`{"Descripcion": "Mastocitoma grado II"`)
	bodyReader := bytes.NewBuffer(jsonBody)

	req := httptest.NewRequest(http.MethodPost, "/api/descripciones?protocolo=CAN-001", bodyReader)
	rec := httptest.NewRecorder()

	handler.APICreateDescripcion(rec, req)

	// Afirmaciones
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error procesando el JSON")
}
func TestAPIGetDescripciones_parametrosValidos_retornaStatusOK(t *testing.T) {
	mockRepo := new(mocks.MockDescripcionMicroscopicaRepository)
	service := application.NewDescripcionMicroscopicaService(mockRepo)
	handler := ui.NewDescripcionMicroscopicaHandler(service)

	esperado := []domain.Descripcion_microscopicas{
		{Descripcion: "Tejido inflamado"},
	}
	mockRepo.On("GetDescripcionMicroscopicaByProtocolo", mock.Anything, "CAN-001").Return(esperado, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/descripciones?protocolo=CAN-001", nil)
	rec := httptest.NewRecorder()

	handler.APIGetDescripciones(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Tejido inflamado")
	mockRepo.AssertExpectations(t)
}

func TestAPIGetDescripciones_faltaProtocolo_retornaBadRequest(t *testing.T) {
	mockRepo := new(mocks.MockDescripcionMicroscopicaRepository)
	service := application.NewDescripcionMicroscopicaService(mockRepo)
	handler := ui.NewDescripcionMicroscopicaHandler(service)

	// Request sin el query param de protocolo
	req := httptest.NewRequest(http.MethodGet, "/api/descripciones", nil)
	rec := httptest.NewRecorder()

	handler.APIGetDescripciones(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Falta el parámetro protocolo")
}
