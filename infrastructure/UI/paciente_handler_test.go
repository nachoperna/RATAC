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

func TestListPacientes_ejecucionExitosa_renderizaTemplate(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)
	handler := ui.NewPacienteHandler(service)

	mockRepo.On("ListPacientes", mock.Anything).Return([]domain.Paciente{}, int16(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/pacientes", nil)
	rec := httptest.NewRecorder()

	handler.ListPacientes(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestListPacientesBy_conParametro_llamaAlServicio(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)
	handler := ui.NewPacienteHandler(service)

	esperado := []domain.Paciente{{NombrePaciente: "Fido"}}
	mockRepo.On("GetPacienteByNombre", mock.Anything, "Fido", int8(0)).Return(esperado, int16(1), nil)

	req := httptest.NewRequest(http.MethodGet, "/pacientes/nombre?paciente=Fido", nil)
	rec := httptest.NewRecorder()

	handler.ListPacientesBy(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestListPacientesByFiltro_jsonValido_llamaAlServicio(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)
	handler := ui.NewPacienteHandler(service)

	esperado := []domain.Paciente{{Protocolo: "FEL-002"}}
	mockRepo.On("GetPacienteByFiltro", mock.Anything, mock.AnythingOfType("[]domain.Filtro"), int8(0)).Return(esperado, int16(1), nil)

	jsonBody := []byte(`{"filtros": [{"campo": "Especie", "operador": "=", "valores": ["Felino"]}], "offset": 0}`)
	req := httptest.NewRequest(http.MethodPost, "/pacientes/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ListPacientesByFiltro(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestListPacientesByFiltro_jsonInvalido_retornaBadRequest(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)
	handler := ui.NewPacienteHandler(service)

	jsonBody := []byte(`{"filtros": bad_json}`) // JSON roto
	req := httptest.NewRequest(http.MethodPost, "/pacientes/", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()

	handler.ListPacientesByFiltro(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRepo.AssertNotCalled(t, "GetPacienteByFiltro")
}

func TestAPIPacientes_ejecucionExitosa_retornaJSON(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)
	handler := ui.NewPacienteHandler(service)

	// Le pasamos solo 2 pacientes para probar que el límite seguro (limite = 5) funciona bien
	esperado := []domain.Paciente{
		{Protocolo: "CAN-001"},
		{Protocolo: "CAN-002"},
	}
	mockRepo.On("ListPacientes", mock.Anything).Return(esperado, int16(len(esperado)), nil)

	req := httptest.NewRequest(http.MethodGet, "/apipacientes", nil)
	rec := httptest.NewRecorder()

	handler.APIPacientes(rec, req)

	// Afirmaciones
	assert.Equal(t, http.StatusOK, rec.Code)

	// Validamos que el JSON devuelto contenga los protocolos esperados
	assert.Contains(t, rec.Body.String(), "CAN-001")
	assert.Contains(t, rec.Body.String(), "CAN-002")

	mockRepo.AssertExpectations(t)
}

func TestShowFullPaciente_protocoloValido_renderizaTemplate(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)
	handler := ui.NewPacienteHandler(service)

	protocolo := "CAN-010"
	esperado := &domain.Paciente{Protocolo: protocolo}
	mockRepo.On("GetAllFromPaciente", mock.Anything, protocolo).Return(esperado, nil)

	req := httptest.NewRequest(http.MethodGet, "/paciente/protocolo/CAN-010", nil)
	// Feature de Go 1.22: Inyectar el PathValue directamente en el request mockeado
	req.SetPathValue("protocolo", protocolo)

	rec := httptest.NewRecorder()

	handler.ShowFullPaciente(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}
