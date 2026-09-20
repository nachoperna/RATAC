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

-- name: GetEmail :one
SELECT u.email
FROM Sesiones s JOIN Usuarios u ON s.id_usuario = u.id
WHERE s.token = $1;
