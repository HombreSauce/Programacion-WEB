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
    id_turno SERIAL  NOT NULL,
    id_medico int  NOT NULL,
    id_paciente int  NOT NULL,
    estado varchar(10)  NOT NULL,
    fecha date  NOT NULL,
    hora time  NOT NULL,
    CONSTRAINT turno_PK PRIMARY KEY (id_turno)
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
    ON DELETE CASCADE
;

-- Reference: atiende_por_obra_social (table: atiende_por)
ALTER TABLE atiende_por ADD CONSTRAINT atiende_por_obra_social_FK
    FOREIGN KEY (obra_social_nombre)
    REFERENCES obra_social (nombre)  
    ON DELETE CASCADE
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

-- Restrcciones y reglas adicionales para turnos

-- Estado con CHECK (proyecto chico, flexible)
ALTER TABLE turnos
  ALTER COLUMN estado SET NOT NULL,
  ADD CONSTRAINT turnos_estado_chk
  CHECK (estado IN ('programado','atendido','cancelado'));

-- Default al crear turno
ALTER TABLE turnos
  ALTER COLUMN estado SET DEFAULT 'programado';

-- Evitar solapado por médico (solo turnos activos)
CREATE UNIQUE INDEX IF NOT EXISTS ux_turnos_medico_fecha_hora_activos
  ON turnos (id_medico, fecha, hora)
  WHERE estado IN ('programado');

-- No permitir turnos en el pasado
 ALTER TABLE turnos
   ADD CONSTRAINT turnos_no_pasados
   CHECK ( (fecha > CURRENT_DATE) OR (fecha = CURRENT_DATE AND hora > CURRENT_TIME) );

-- No permitir que un medico sea su propio paciente
ALTER TABLE turnos
  ADD CONSTRAINT turnos_medico_distinto_de_paciente
  CHECK (id_medico <> id_paciente);

-- Inserta si no existe (seguro para ejecutar varias veces)
INSERT INTO obra_social (nombre)
VALUES
  ('PAMI - INSSJP'),
  ('IOMA'),
  ('OSDE'),
  ('OSECAC - Empleados de Comercio'),
  ('OSPRERA - Trabajadores Rurales'),
  ('OSPe - Petroleros'),
  ('OSUOMRA - Metalúrgicos'),
  ('OSPJN - Judiciales Nación'),
  ('OSMATA - Madereros'),
  ('OSDEPYM - Dirigentes de Empresas')
ON CONFLICT (nombre) DO NOTHING;

