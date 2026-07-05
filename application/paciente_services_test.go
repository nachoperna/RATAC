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

func TestListPacientes_ejecucionExitosa_retornaLista(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)

	esperado := []domain.Paciente{
		{Protocolo: "CAN-001", NombrePaciente: "Fido"},
		{Protocolo: "FEL-002", NombrePaciente: "Michi"},
	}
	var cant int16 = 2

	mockRepo.On("ListPacientes", mock.Anything, 0).Return(esperado, cant, nil)

	resultado, total, err := service.ListPacientes(context.Background(), 0)

	assert.NoError(t, err)
	assert.Equal(t, esperado, resultado)
	assert.Equal(t, cant, total)
	mockRepo.AssertExpectations(t)
}

func TestListUltimosPacientes_ejecucionExitosa_retornaLista(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)

	esperado := []domain.Paciente{{Protocolo: "CAN-999"}}

	mockRepo.On("ListUltimosPacientes", mock.Anything).Return(esperado, []bool{}, nil)

	resultado, _, err := service.ListUltimosPacientes(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, esperado, resultado)
	mockRepo.AssertExpectations(t)
}

func TestCountPacientes_ejecucionExitosa_retornaTotal(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)

	var esperado int64 = 150

	mockRepo.On("CountPacientes", mock.Anything).Return(esperado, nil)

	resultado, err := service.CountPacientes(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, esperado, resultado)
	mockRepo.AssertExpectations(t)
}

func TestGetPacienteByNombre_parametrosValidos_retornaPacientesYTotal(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)

	nombreBuscado := "Luna"
	var offset int8 = 0
	esperadoLista := []domain.Paciente{{NombrePaciente: "Luna"}}
	var esperadoTotal int16 = 1

	// Notar que acá retornamos 3 valores: la lista, el total y el error (nil)
	mockRepo.On("GetPacienteByNombre", mock.Anything, nombreBuscado, offset).Return(esperadoLista, esperadoTotal, nil)

	resultadoLista, resultadoTotal, err := service.GetPacienteByNombre(context.Background(), nombreBuscado, offset)

	assert.NoError(t, err)
	assert.Equal(t, esperadoLista, resultadoLista)
	assert.Equal(t, esperadoTotal, resultadoTotal)
	mockRepo.AssertExpectations(t)
}

func TestGetPacienteByFiltro_parametrosValidos_retornaPacientesYTotal(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)

	filtros := []domain.Filtro{
		{Campo: "Especie", Operador: "=", Valores: []string{"Canino"}},
	}
	var offset int8 = 10
	esperadoLista := []domain.Paciente{{Protocolo: "CAN-010"}}
	var esperadoTotal int16 = 45

	mockRepo.On("GetPacienteByFiltro", mock.Anything, filtros, offset).Return(esperadoLista, esperadoTotal, nil)

	resultadoLista, resultadoTotal, err := service.GetPacienteByFiltro(context.Background(), filtros, offset)

	assert.NoError(t, err)
	assert.Equal(t, esperadoLista, resultadoLista)
	assert.Equal(t, esperadoTotal, resultadoTotal)
	mockRepo.AssertExpectations(t)
}

func TestGetAllFromPaciente_protocoloValido_retornaPacienteCompleto(t *testing.T) {
	mockRepo := new(mocks.MockPacienteRepository)
	service := application.NewPacienteService(mockRepo)

	protocolo := "CAN-005"
	especie := "Canino"
	esperado := &domain.Paciente{
		Protocolo: protocolo,
		Especie:   &especie,
	}

	mockRepo.On("GetAllFromPaciente", mock.Anything, protocolo).Return(esperado, nil)

	resultado, err := service.GetAllFromPaciente(context.Background(), protocolo)

	assert.NoError(t, err)
	assert.Equal(t, esperado, resultado)
	mockRepo.AssertExpectations(t)
}
