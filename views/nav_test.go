package views_test

import (
	"RATAC/domain"
	"RATAC/views"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func renderNav(t *testing.T, u *domain.Usuario) string {
	t.Helper()
	var sb strings.Builder
	err := views.Nav(u).Render(context.Background(), &sb)
	assert.NoError(t, err)
	return sb.String()
}

// El bug que cierra el nav unico: los estaticos ofrecian "Cargar Diagnostico"
// y "Mi Panel" a cualquier anonimo.
func TestNav_anonimo_soloMuestraIngresar(t *testing.T) {
	html := renderNav(t, domain.UsuarioPublico())

	assert.Contains(t, html, "/login")
	assert.Contains(t, html, "Ingresar")
	assert.NotContains(t, html, "/diagnosticos/alta")
	assert.NotContains(t, html, "/admin/panel")
	assert.NotContains(t, html, "/logout")
}

func TestNav_autenticado_muestraPanelYSalir(t *testing.T) {
	for _, rol := range []domain.Rol{domain.RolAdmin, domain.RolLaboratorio} {
		html := renderNav(t, &domain.Usuario{ID: 1, Usuario: "u", Rol: rol})

		assert.Contains(t, html, "/diagnosticos/alta", "rol %q", rol)
		assert.Contains(t, html, "/admin/panel", "rol %q", rol)
		assert.Contains(t, html, "/logout", "rol %q", rol)
		assert.NotContains(t, html, "Ingresar", "rol %q", rol)
	}
}

// Nav se alimenta de UsuarioDeCtx, que ante un context vacio da el publico.
func TestNav_contextSinUsuario_rindeNavPublico(t *testing.T) {
	html := renderNav(t, domain.UsuarioDeCtx(context.Background()))
	assert.Contains(t, html, "Ingresar")
	assert.NotContains(t, html, "/admin/panel")
}
