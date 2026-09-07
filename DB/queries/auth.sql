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

-- name: NuevaSesion :one
INSERT INTO Sesiones (
      token, id_usuario, expiracion, fecha_creacion
) VALUES ( $1, $2, NOW() + INTERVAL '2 minutes', default )
RETURNING *;

-- name: DeleteSesionesViejas :exec
DELETE FROM Sesiones WHERE id_usuario = $1 AND expiracion < NOW();

-- name: DeleteSesion :exec
DELETE FROM Sesiones WHERE token = $1;

-- name: SesionActiva :one
SELECT EXISTS (
      SELECT 1 FROM Sesiones WHERE token = $1 AND expiracion >= NOW()
);
