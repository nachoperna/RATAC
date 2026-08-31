-- Table: Pacientes
CREATE TABLE Pacientes (
    id serial NOT NULL,
    Protocolo varchar(50) NOT NULL UNIQUE,
    Fecha date  NOT NULL,
    Solicitante varchar(100) NOT NULL,
    Tecnica varchar(15) NOT NULL,
    Familia varchar(100) NULL,
    Especie varchar(6) NULL,
    Raza varchar(100),
    Edad smallint,
    Paciente varchar(100) NOT NULL,
    Antecedentes text NULL,
    Descripcion_macroscopica text,
    Referencias_mastocitomas boolean NOT NULL,
    CONSTRAINT Pacientes_pk PRIMARY KEY (id,Protocolo)
);

-- Table: Descripciones_microscopicas
CREATE TABLE Descripciones_microscopicas (
    Descripcion text NOT NULL,
    Diagnostico text NULL,
    Pacientes_Protocolo varchar(50) NOT NULL,
    CONSTRAINT Descripciones_microscopicas_pk PRIMARY KEY (Descripcion,Pacientes_Protocolo)
);

-- Table: Imagenes
CREATE TABLE Imagenes (
    Ruta text NOT NULL,
    Descripciones_microscopicas_Descripcion text NOT NULL,
    Descripciones_microscopicas_Pacientes_Protocolo varchar(50) NOT NULL,
    CONSTRAINT Imagenes_pk PRIMARY KEY (Ruta)
);

-- Table: Grado_oncologico
CREATE TABLE Grado_oncologico (
    id serial NOT NULL,
    Caracteristica varchar(100) NOT NULL,
    Muestra_analizada varchar(100) NULL,
    Puntaje smallint NOT NULL,
    Descripciones_microscopicas_Descripcion text NOT NULL,
    Descripciones_microscopicas_Pacientes_Protocolo varchar(50) NOT NULL,
    CONSTRAINT Grado_oncologico_pk PRIMARY KEY (id,Descripciones_microscopicas_Pacientes_Protocolo,Descripciones_microscopicas_Descripcion)
);

-- Table: usuarios
CREATE TABLE usuarios (
    id            SERIAL PRIMARY KEY,
    usuario       VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(60)  NOT NULL,
    nombre        VARCHAR(100) NOT NULL,
    rol           VARCHAR(20)  NOT NULL DEFAULT 'laboratorio'
                  CHECK (rol IN ('admin','laboratorio')),
    activo        BOOLEAN      NOT NULL DEFAULT TRUE,
    creado_en     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE historial_acciones (
--     id serial PRIMARY KEY,
--     usuario varchar(50) NOT NULL DEFAULT 'sistema',
--     accion varchar(255) NOT NULL,
--     tabla_afectada varchar(50) NOT NULL,
--     fecha timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
--    datos_anteriores text,
--    datos_nuevos text
--);


-- Funciones y triggers para registrar el historial de acciones

-- CREATE OR REPLACE FUNCTION registrar_historial()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     IF (TG_OP = 'DELETE') THEN
--         -- Si borramos, solo guardamos lo que se eliminó (OLD)
--         INSERT INTO historial_acciones (accion, tabla_afectada, datos_viejos)
--         VALUES (TG_OP, TG_TABLE_NAME, row_to_json(OLD)::TEXT);
--         RETURN OLD;
--
--     ELSIF (TG_OP = 'UPDATE') THEN
--         -- Si actualizamos, guardamos cómo estaba (OLD) y cómo quedó (NEW)
--         INSERT INTO historial_acciones (accion, tabla_afectada, datos_viejos, datos_nuevos)
--         VALUES (TG_OP, TG_TABLE_NAME, row_to_json(OLD)::TEXT, row_to_json(NEW)::TEXT);
--         RETURN NEW;
--
--     ELSIF (TG_OP = 'INSERT') THEN
--         -- Si insertamos, solo guardamos lo nuevo (NEW)
--         INSERT INTO historial_acciones (accion, tabla_afectada, datos_nuevos)
--         VALUES (TG_OP, TG_TABLE_NAME, row_to_json(NEW)::TEXT);
--         RETURN NEW;
--
--     END IF;
--
--     RETURN NULL;
-- END;
-- $$ LANGUAGE plpgsql;
--
-- CREATE TRIGGER trigger_historial_diagnosticos
-- AFTER INSERT OR UPDATE OR DELETE ON Pacientes, descripciones_microscopicas, Imagenes, Grado_oncologico
-- FOR EACH ROW
-- EXECUTE FUNCTION registrar_historial();
