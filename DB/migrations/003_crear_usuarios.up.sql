CREATE TABLE IF NOT EXISTS usuarios (
    id            SERIAL PRIMARY KEY,
    usuario       VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(60)  NOT NULL,
    nombre        VARCHAR(100) NOT NULL,
    rol           VARCHAR(20)  NOT NULL DEFAULT 'laboratorio'
                  CHECK (rol IN ('admin','laboratorio')),
    activo        BOOLEAN      NOT NULL DEFAULT TRUE,
    creado_en     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Admin inicial. El hash bcrypt se genero una unica vez fuera de la migracion
-- para que esta sea determinista; la contrasena se comunica aparte y debe
-- cambiarse tras el primer login.
INSERT INTO usuarios (usuario, password_hash, nombre, rol)
VALUES ('admin', '$2a$10$wG9/z1UV22rUkYtGJ.Qfvu8fBA32G3zyd3WM0Y0G4RKE/AtAqAICW', 'Administrador RATAC', 'admin')
ON CONFLICT (usuario) DO NOTHING;
