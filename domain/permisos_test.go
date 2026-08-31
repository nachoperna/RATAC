package domain_test

import (
	"RATAC/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPuede_matrizRolPermiso(t *testing.T) {
	todos := []domain.Permiso{
		domain.PermVerPanel,
		domain.PermCargarDiag,
		domain.PermBorrarDiag,
		domain.PermVerDiagnosticos,
		domain.PermGestionUsuarios,
	}

	esperado := map[domain.Rol]map[domain.Permiso]bool{
		domain.RolPublico: {},
		domain.RolLaboratorio: {
			domain.PermVerPanel:        true,
			domain.PermCargarDiag:      true,
			domain.PermBorrarDiag:      true,
			domain.PermVerDiagnosticos: true,
		},
		domain.RolAdmin: {
			domain.PermVerPanel:        true,
			domain.PermCargarDiag:      true,
			domain.PermBorrarDiag:      true,
			domain.PermVerDiagnosticos: true,
			domain.PermGestionUsuarios: true,
		},
	}

	for rol, permisos := range esperado {
		u := &domain.Usuario{Rol: rol}
		for _, p := range todos {
			assert.Equal(t, permisos[p], domain.Puede(u, p),
				"rol %q y permiso %q", rol, p)
		}
	}
}

func TestPuede_laboratorioNoGestionaUsuarios(t *testing.T) {
	lab := &domain.Usuario{Rol: domain.RolLaboratorio}
	assert.False(t, domain.Puede(lab, domain.PermGestionUsuarios))
	assert.True(t, domain.Puede(lab, domain.PermCargarDiag))
}

func TestPuede_publicoNoTieneNingunPermiso(t *testing.T) {
	publico := domain.UsuarioPublico()
	for _, p := range []domain.Permiso{
		domain.PermVerPanel, domain.PermCargarDiag, domain.PermBorrarDiag,
		domain.PermVerDiagnosticos, domain.PermGestionUsuarios,
	} {
		assert.False(t, domain.Puede(publico, p), "permiso %q", p)
	}
}

func TestPuede_usuarioNil_retornaFalse(t *testing.T) {
	assert.False(t, domain.Puede(nil, domain.PermVerPanel))
}

func TestAutenticado(t *testing.T) {
	assert.False(t, domain.UsuarioPublico().Autenticado())
	assert.False(t, (*domain.Usuario)(nil).Autenticado())
	assert.False(t, (&domain.Usuario{}).Autenticado())
	assert.True(t, (&domain.Usuario{Rol: domain.RolAdmin}).Autenticado())
	assert.True(t, (&domain.Usuario{Rol: domain.RolLaboratorio}).Autenticado())
}
