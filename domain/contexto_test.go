package domain_test

import (
	"RATAC/domain"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsuarioDeCtx_contextVacio_devuelvePublico(t *testing.T) {
	u := domain.UsuarioDeCtx(context.Background())
	assert.NotNil(t, u, "nunca debe devolver nil")
	assert.Equal(t, domain.RolPublico, u.Rol)
	assert.False(t, u.Autenticado())
}

func TestUsuarioDeCtx_conUsuario_loRecupera(t *testing.T) {
	esperado := &domain.Usuario{ID: 7, Usuario: "admin", Rol: domain.RolAdmin}
	ctx := domain.ConUsuarioCtx(context.Background(), esperado)

	u := domain.UsuarioDeCtx(ctx)
	assert.Equal(t, esperado, u)
	assert.True(t, u.Autenticado())
}
