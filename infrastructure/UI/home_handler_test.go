package ui_test

import (
	"RATAC/application"
	ui "RATAC/infrastructure/UI"
	"RATAC/mocks"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShowHome_rutaInvalida_retornaArchivoError(t *testing.T) {
	// Guardamos dónde estamos y subimos a la raíz (RATAC/) para que encuentre el HTML
	dirOriginal, _ := os.Getwd()
	os.Chdir("../..")
	defer os.Chdir(dirOriginal) // Esto asegura que al terminar el test volvemos al lugar original

	mockPacRepo := new(mocks.MockPacienteRepository)
	mockDiagRepo := new(mocks.MockDiagnosticoRepository)
	mockDescRepo := new(mocks.MockDescripcionMicroscopicaRepository)

	srvPac := application.NewPacienteService(mockPacRepo)
	srvDiag := application.NewDiagnosticoService(mockDiagRepo)
	srvDesc := application.NewDescripcionMicroscopicaService(mockDescRepo)

	handler := ui.NewHomeHandler(srvPac, srvDesc, srvDiag)

	req := httptest.NewRequest(http.MethodGet, "/ruta-que-no-existe", nil)
	rec := httptest.NewRecorder()

	handler.ShowHome(rec, req)

	mockPacRepo.AssertNotCalled(t, "CountPacientes")
}

func TestShowHome_rutaValida_ejecutaTemplate(t *testing.T) {
	// Guardamos dónde estamos y subimos a la raíz (RATAC/) para que encuentre el HTML
	dirOriginal, _ := os.Getwd()
	os.Chdir("../..")
	defer os.Chdir(dirOriginal)

	mockPacRepo := new(mocks.MockPacienteRepository)
	mockDiagRepo := new(mocks.MockDiagnosticoRepository)
	mockDescRepo := new(mocks.MockDescripcionMicroscopicaRepository)

	srvPac := application.NewPacienteService(mockPacRepo)
	srvDiag := application.NewDiagnosticoService(mockDiagRepo)
	srvDesc := application.NewDescripcionMicroscopicaService(mockDescRepo)

	handler := ui.NewHomeHandler(srvPac, srvDesc, srvDiag)

	mockPacRepo.On("CountPacientes", mock.Anything).Return(int64(10), nil)
	mockDiagRepo.On("CountImagenes", mock.Anything).Return(int64(5), nil)
	mockDescRepo.On("CountDiagnosticos", mock.Anything).Return(int64(8), nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ShowHome(rec, req)

	// Ahora sí, debería dar 200 OK
	assert.Equal(t, http.StatusOK, rec.Code)
	mockPacRepo.AssertExpectations(t)
	mockDiagRepo.AssertExpectations(t)
	mockDescRepo.AssertExpectations(t)
}
