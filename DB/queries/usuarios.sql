-- name: GetUsuarioByNombre :one
SELECT * FROM usuarios
WHERE usuario = $1 AND activo = TRUE
LIMIT 1;

-- name: GetUsuarioByID :one
SELECT * FROM usuarios
WHERE id = $1 AND activo = TRUE
LIMIT 1;

-- name: CreateUsuario :one
INSERT INTO usuarios (
    usuario, password_hash, nombre, rol
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdatePasswordUsuario :exec
UPDATE usuarios
SET password_hash = $2
WHERE id = $1;
