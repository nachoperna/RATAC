-- Vuelve a la version de 001: el usuario queda siempre en el DEFAULT 'sistema'.
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
