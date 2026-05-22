-- 1. Crear la tabla del historial
CREATE TABLE IF NOT EXISTS historial_acciones (
    id SERIAL PRIMARY KEY,
    accion VARCHAR(50) NOT NULL,
    tabla_afectada VARCHAR(50) NOT NULL,
    datos_viejos TEXT,  -- Estado anterior (OLD)
    datos_nuevos TEXT,  -- Estado nuevo (NEW)
    fecha TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    usuario VARCHAR(50) DEFAULT 'sistema'
);

-- 2. Crear la función dinámica
CREATE OR REPLACE FUNCTION registrar_historial()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO historial_acciones (accion, tabla_afectada, datos_viejos)
        VALUES (TG_OP, TG_TABLE_NAME, row_to_json(OLD)::TEXT);
        RETURN OLD;
        
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO historial_acciones (accion, tabla_afectada, datos_viejos, datos_nuevos)
        VALUES (TG_OP, TG_TABLE_NAME, row_to_json(OLD)::TEXT, row_to_json(NEW)::TEXT);
        RETURN NEW;
        
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO historial_acciones (accion, tabla_afectada, datos_nuevos)
        VALUES (TG_OP, TG_TABLE_NAME, row_to_json(NEW)::TEXT);
        RETURN NEW;
        
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- 3. Enganchar los Triggers de forma individual para cada tabla
CREATE TRIGGER trigger_historial_pacientes
AFTER INSERT OR UPDATE OR DELETE ON Pacientes
FOR EACH ROW EXECUTE FUNCTION registrar_historial();

CREATE TRIGGER trigger_historial_descripciones
AFTER INSERT OR UPDATE OR DELETE ON Descripciones_microscopicas
FOR EACH ROW EXECUTE FUNCTION registrar_historial();

CREATE TRIGGER trigger_historial_imagenes
AFTER INSERT OR UPDATE OR DELETE ON Imagenes
FOR EACH ROW EXECUTE FUNCTION registrar_historial();

CREATE TRIGGER trigger_historial_grado
AFTER INSERT OR UPDATE OR DELETE ON Grado_oncologico
FOR EACH ROW EXECUTE FUNCTION registrar_historial();