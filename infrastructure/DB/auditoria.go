package db

import (
	"RATAC/domain"
	"context"
	"database/sql"
)

// fijarUsuarioAuditoria deja el usuario de la aplicacion en una variable de
// sesion que lee el trigger registrar_historial(); sin esto el historial cae
// siempre en el DEFAULT 'sistema'.
//
// Se usa set_config con parametro y no SET LOCAL porque SET LOCAL no admite
// placeholders y obligaria a interpolar el string en el SQL. El TRUE final lo
// hace local a la transaccion, evitando fugas entre conexiones del pool.
func fijarUsuarioAuditoria(ctx context.Context, tx *sql.Tx) error {
	usuario := domain.UsuarioDeCtx(ctx)
	nombre := usuario.Usuario
	if nombre == "" {
		nombre = "sistema"
	}
	_, err := tx.ExecContext(ctx, "SELECT set_config('ratac.usuario', $1, TRUE)", nombre)
	return err
}
