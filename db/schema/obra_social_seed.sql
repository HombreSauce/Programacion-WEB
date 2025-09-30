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