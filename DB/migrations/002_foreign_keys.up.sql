ALTER TABLE Pacientes ADD CONSTRAINT Pacientes_Protocolo_unique UNIQUE (Protocolo);

-- foreign keys
-- Reference: Paciente_Laboratorio (table: Pacientes)
ALTER TABLE Pacientes ADD CONSTRAINT Paciente_Laboratorio
      FOREIGN KEY (Email_Lab)
      REFERENCES Usuarios (email) 
      NOT DEFERRABLE 
      INITIALLY IMMEDIATE
;

-- Reference: Descripciones_microscopicas_Pacientes (table: Descripciones_microscopicas)
ALTER TABLE Descripciones_microscopicas ADD CONSTRAINT Descripciones_microscopicas_Pacientes
      FOREIGN KEY (Pacientes_Protocolo)
      REFERENCES Pacientes (Protocolo) 
      ON DELETE CASCADE
      NOT DEFERRABLE 
      INITIALLY IMMEDIATE
;

-- Reference: Imagenes_Descripciones_microscopicas (table: Imagenes)
ALTER TABLE Imagenes ADD CONSTRAINT Imagenes_Descripciones_microscopicas
      FOREIGN KEY (Descripciones_microscopicas_Descripcion, Descripciones_microscopicas_Pacientes_Protocolo)
      REFERENCES Descripciones_microscopicas (Descripcion, Pacientes_Protocolo)  
      ON DELETE CASCADE
      NOT DEFERRABLE 
      INITIALLY IMMEDIATE
;

-- Reference: Grado_oncologico_Descripciones_microscopicas (table: Grado_oncologico)
ALTER TABLE Grado_oncologico ADD CONSTRAINT Grado_oncologico_Descripciones_microscopicas
      FOREIGN KEY (Descripciones_microscopicas_Descripcion, Descripciones_microscopicas_Pacientes_Protocolo)
      REFERENCES Descripciones_microscopicas (Descripcion, Pacientes_Protocolo)  
      ON DELETE CASCADE
      NOT DEFERRABLE 
      INITIALLY IMMEDIATE
;
