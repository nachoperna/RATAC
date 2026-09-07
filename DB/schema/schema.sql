-- Table: Pacientes
CREATE TABLE Pacientes (
      id serial NOT NULL,
      Protocolo varchar(50) NOT NULL UNIQUE,
      Fecha date NOT NULL,
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

CREATE TYPE roles AS ENUM ('admin', 'laboratorio');

-- Table: Usuarios
CREATE TABLE Usuarios (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      email VARCHAR(255) UNIQUE NOT NULL,
      contraseña_hash VARCHAR(255) NOT NULL,
      nombre_lab VARCHAR(255) NOT NULL,
      rol roles NOT NULL,
      ciudad_origen VARCHAR(255) NOT NULL,
      activo BOOLEAN,
      fecha_creacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      fecha_modificacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Table: Sesiones
CREATE TABLE Sesiones (
      token VARCHAR(255) PRIMARY KEY,
      id_usuario UUID NOT NULL REFERENCES Usuarios(id) ON DELETE CASCADE,
      expiracion TIMESTAMPTZ NOT NULL,
      fecha_creacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sesiones_token ON Sesiones(token);
