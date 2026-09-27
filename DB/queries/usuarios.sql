-- name: CreateUsuario :one
INSERT INTO Usuarios (
      id, email, nombre_lab, contraseña_hash, rol, ciudad_origen, activo, fecha_creacion, fecha_modificacion
) VALUES (default, $1, $2, $3, $4, $5, true, default, default)
RETURNING *;

-- name: GetUsuario :one
SELECT * FROM Usuarios WHERE email = $1;

-- name: UpdateContraseña :exec
UPDATE Usuarios SET contraseña_hash = $2 WHERE email = $1;

-- name: Login :one
SELECT * FROM Usuarios WHERE email = $1;

-- name: UsuarioLogueado :exec
UPDATE Usuarios SET activo = true WHERE email = $1;

-- name: CambiarRol :exec
UPDATE Usuarios set rol = $2::roles WHERE email = $1;

-- name: EliminarUsuario :one
DELETE FROM Usuarios WHERE email = $1
RETURNING nombre_lab;
