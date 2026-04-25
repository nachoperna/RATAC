-- Table: Pacientes
CREATE TABLE Pacientes (
    id serial NOT NULL,
    Protocolo varchar(20) NOT NULL,
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
    CONSTRAINT Pacientes_pk PRIMARY KEY (Protocolo)
);

-- Table: Descripciones_microscopicas
CREATE TABLE Descripciones_microscopicas (
    Descripcion text NOT NULL,
    Diagnostico text NULL,
    Pacientes_Protocolo varchar(20) NOT NULL,
    CONSTRAINT Descripciones_microscopicas_pk PRIMARY KEY (Descripcion,Pacientes_Protocolo)
);

-- Table: Imagenes
CREATE TABLE Imagenes (
    Ruta text NOT NULL,
    Descripciones_microscopicas_Descripcion text NOT NULL,
    Descripciones_microscopicas_Pacientes_Protocolo varchar(20) NOT NULL,
    CONSTRAINT Imagenes_pk PRIMARY KEY (Ruta)
);

-- Table: Grado_oncologico
CREATE TABLE Grado_oncologico (
    Caracteristica varchar(100) NOT NULL,
    Muestra_analizada varchar(100) NULL,
    Puntaje smallint NOT NULL,
    Total smallint NOT NULL,
    Descripciones_microscopicas_Descripcion text NOT NULL,
    Descripciones_microscopicas_Pacientes_Protocolo varchar(20) NOT NULL,
    CONSTRAINT Grado_oncologico_pk PRIMARY KEY (Descripciones_microscopicas_Pacientes_Protocolo,Descripciones_microscopicas_Descripcion)
);

-- foreign keys
-- Reference: Descripciones_microscopicas_Pacientes (table: Descripciones_microscopicas)
ALTER TABLE Descripciones_microscopicas ADD CONSTRAINT Descripciones_microscopicas_Pacientes
    FOREIGN KEY (Pacientes_Protocolo)
    REFERENCES Pacientes (Protocolo) 
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Imagenes_Descripciones_microscopicas (table: Imagenes)
ALTER TABLE Imagenes ADD CONSTRAINT Imagenes_Descripciones_microscopicas
    FOREIGN KEY (Descripciones_microscopicas_Descripcion, Descripciones_microscopicas_Pacientes_Protocolo)
    REFERENCES Descripciones_microscopicas (Descripcion, Pacientes_Protocolo)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Grado_oncologico_Descripciones_microscopicas (table: Grado_oncologico)
ALTER TABLE Grado_oncologico ADD CONSTRAINT Grado_oncologico_Descripciones_microscopicas
    FOREIGN KEY (Descripciones_microscopicas_Descripcion, Descripciones_microscopicas_Pacientes_Protocolo)
    REFERENCES Descripciones_microscopicas (Descripcion, Pacientes_Protocolo)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;
