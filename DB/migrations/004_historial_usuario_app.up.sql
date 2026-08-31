-- El trigger plpgsql no conoce al usuario de la aplicacion: se le pasa por
-- variable de sesion (ratac.usuario), que la app fija por transaccion con
-- set_config(..., TRUE).
CREATE OR REPLACE FUNCTION registrar_historial()
RETURNS TRIGGER AS $$
DECLARE
    usuario_app VARCHAR(50);
BEGIN
    -- El TRUE de current_setting es obligatorio: sin el, la funcion falla
    -- cuando la variable no esta seteada (p. ej. escrituras desde ProcesadoJsons/).
    usuario_app := COALESCE(NULLIF(current_setting('ratac.usuario', TRUE), ''), 'sistema');

    IF (TG_OP = 'DELETE') THEN
        INSERT INTO historial_acciones (accion, tabla_afectada, datos_viejos, usuario)
        VALUES (TG_OP, TG_TABLE_NAME, row_to_json(OLD)::TEXT, usuario_app);
        RETURN OLD;

    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO historial_acciones (accion, tabla_afectada, datos_viejos, datos_nuevos, usuario)
        VALUES (TG_OP, TG_TABLE_NAME, row_to_json(OLD)::TEXT, row_to_json(NEW)::TEXT, usuario_app);
        RETURN NEW;

    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO historial_acciones (accion, tabla_afectada, datos_nuevos, usuario)
        VALUES (TG_OP, TG_TABLE_NAME, row_to_json(NEW)::TEXT, usuario_app);
        RETURN NEW;

    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
