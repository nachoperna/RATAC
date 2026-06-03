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

func TestAPIGetDiagnostico_parametrosFaltantes_retornaBadRequest(t *testing.T) {
	mockRepo := new(mocks.MockDiagnosticoRepository)
	service := application.NewDiagnosticoService(mockRepo)
	handler := ui.NewDiagnosticoHandler(service)

	// Escenario: request GET sin query params
	req := httptest.NewRequest(http.MethodGet, "/api/diagnostico", nil)
	rec := httptest.NewRecorder()

	handler.APIGetDiagnostico(rec, req)

	// Resultado: Código 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Faltan parámetros en la URL")
}

func TestAPIGetDiagnostico_parametrosValidos_retornaStatusOK(t *testing.T) {
	mockRepo := new(mocks.MockDiagnosticoRepository)
	service := application.NewDiagnosticoService(mockRepo)
	handler := ui.NewDiagnosticoHandler(service)

	// Escenario: seteamos el mock para que devuelva un diagnóstico correctamente
	descTexto := "Diagnostico positivo"
	esperado := &domain.Diagnostico{
		Descripcion: &descTexto,
		Imagenes:    []string{"img1.png"},
	}
	// El handler le pasa la data al service, y el service al repo
	mockRepo.On("GetDiagnostico", mock.Anything, "1234", "desc_test").Return(esperado, nil)

	// Armamos la request con los parámetros correctos
	req := httptest.NewRequest(http.MethodGet, "/api/diagnostico?protocolo=1234&descripcion=desc_test", nil)
	rec := httptest.NewRecorder()

	handler.APIGetDiagnostico(rec, req)

	// Resultado: Código 200 OK y el JSON correspondiente
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Diagnostico positivo")
	assert.Contains(t, rec.Body.String(), "img1.png")
	mockRepo.AssertExpectations(t)
}

func TestAPICreateDiagnostico_parametrosValidos_retornaStatusCreated(t *testing.T) {
	mockRepo := new(mocks.MockDiagnosticoRepository)
	service := application.NewDiagnosticoService(mockRepo)
	handler := ui.NewDiagnosticoHandler(service)

	mockRepo.On("CreateDiagnostico", mock.Anything, "1234", "desc_test", mock.AnythingOfType("domain.Diagnostico")).Return(nil)

	jsonBody := []byte(`{"Imagenes": ["img1.png", "img2.png"]}`)
	bodyReader := bytes.NewBuffer(jsonBody)

	// Necesita dos query params: protocolo y descripcion
	req := httptest.NewRequest(http.MethodPost, "/api/diagnostico?protocolo=1234&descripcion=desc_test", bodyReader)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.APICreateDiagnostico(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAPIDeleteImagen_parametrosValidos_retornaStatusOK(t *testing.T) {
	mockRepo := new(mocks.MockDiagnosticoRepository)
	service := application.NewDiagnosticoService(mockRepo)
	handler := ui.NewDiagnosticoHandler(service)

	ruta := "IMAGENES/foto.png"
	mockRepo.On("DeleteImagen", mock.Anything, ruta).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/diagnostico/imagen?ruta=IMAGENES/foto.png", nil)
	rec := httptest.NewRecorder()

	handler.APIDeleteImagen(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}
