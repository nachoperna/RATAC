package domain

import "context"

type claveCtx string

const UsuarioCtxKey claveCtx = "usuario"

func ConUsuarioCtx(ctx context.Context, u *Usuario) context.Context {
	return context.WithValue(ctx, UsuarioCtxKey, u)
}

// UsuarioDeCtx nunca devuelve nil: ante un context sin usuario (tests que
// invocan handlers sin pasar por el middleware, renders internos) devuelve
// el usuario publico, que no tiene ningun permiso.
func UsuarioDeCtx(ctx context.Context) *Usuario {
	if ctx == nil {
		return UsuarioPublico()
	}
	if u, ok := ctx.Value(UsuarioCtxKey).(*Usuario); ok && u != nil {
		return u
	}
	return UsuarioPublico()
}
