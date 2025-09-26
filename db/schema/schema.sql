-- Table: atiende_por
CREATE TABLE atiende_por (
    id_medico int  NOT NULL,
    obra_social_nombre varchar(50)  NOT NULL,
    CONSTRAINT atiende_por_PK PRIMARY KEY (id_medico,obra_social_nombre)
);

-- Table: medicos
CREATE TABLE medicos (
    id_medico int  NOT NULL,
    nro_matricula int  NOT NULL,
    especialidad varchar(30)  NOT NULL,
    CONSTRAINT medicos_ak UNIQUE (nro_matricula) NOT DEFERRABLE  INITIALLY IMMEDIATE,
    CONSTRAINT medico_PK PRIMARY KEY (id_medico)
);

-- Table: obra_social
CREATE TABLE obra_social (
    nombre varchar(50)  NOT NULL,
    CONSTRAINT obra_social_PK PRIMARY KEY (nombre)
);

-- Table: pacientes
CREATE TABLE pacientes (
    id_paciente int  NOT NULL,
    obra_social varchar(50)  NULL,
    nro_afiliado varchar(20)  NULL,
    CONSTRAINT paciente_PK PRIMARY KEY (id_paciente)
);

-- Table: turnos
CREATE TABLE turnos (
    id_turno int  NOT NULL,
    id_medico int  NOT NULL,
    id_paciente int  NOT NULL,
    estado varchar(10)  NOT NULL,
    fecha date  NOT NULL,
    hora time  NOT NULL,
    CONSTRAINT turno_PK PRIMARY KEY (ID_turno)
);

-- Table: usuarios
CREATE TABLE usuarios (
    id SERIAL NOT NULL, 
    DNI varchar(9)  NOT NULL,
    nombre varchar(30)  NOT NULL,
    apellido varchar(30)  NOT NULL,
    sexo varchar(10) NOT NULL,
    fecha_nacimiento date  NOT NULL,
    email varchar(50)  NOT NULL,
    telefono varchar(20)  NOT NULL,
    rol char(1)  NOT NULL,
    CONSTRAINT usuarios_PK PRIMARY KEY (ID)
);

-- foreign keys
-- Reference: atiende_por_medicos (table: atiende_por)
ALTER TABLE atiende_por ADD CONSTRAINT atiende_por_medicos_FK
    FOREIGN KEY (id_medico)
    REFERENCES medicos (id_medico)  
;

-- Reference: atiende_por_obra_social (table: atiende_por)
ALTER TABLE atiende_por ADD CONSTRAINT atiende_por_obra_social_FK
    FOREIGN KEY (obra_social_nombre)
    REFERENCES obra_social (nombre)  
;

-- Reference: medicos_usuarios (table: medicos)
ALTER TABLE medicos ADD CONSTRAINT medicos_usuarios_FK
    FOREIGN KEY (id_medico)
    REFERENCES usuarios (ID)  
    ON DELETE CASCADE
;

-- Reference: pacientes_obra_social (table: pacientes)
ALTER TABLE pacientes ADD CONSTRAINT pacientes_obra_social_FK
    FOREIGN KEY (obra_social)
    REFERENCES obra_social (nombre)  
;

-- Reference: pacientes_usuarios (table: pacientes)
ALTER TABLE pacientes ADD CONSTRAINT pacientes_usuarios_FK
    FOREIGN KEY (id_paciente)
    REFERENCES usuarios (ID)  
    ON DELETE CASCADE
;

-- Reference: turnos_medicos (table: turnos)
ALTER TABLE turnos ADD CONSTRAINT turnos_medicos
    FOREIGN KEY (id_medico)
    REFERENCES medicos (id_medico)  
;

-- Reference: turnos_pacientes (table: turnos)
ALTER TABLE turnos ADD CONSTRAINT turnos_pacientes
    FOREIGN KEY (id_paciente)
    REFERENCES pacientes (id_paciente)  
;

-- End of file.

