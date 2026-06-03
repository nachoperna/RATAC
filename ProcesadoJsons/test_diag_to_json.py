import pytest
from diag_to_json import (
    matcheaCategoria, 
    cargarDescMicro, 
    resultados,
    MATERIAL_REMITIDO_ANTECEDENTES,
    DESCRIPCION_MACROSCOPICA,
    DESCRIPCION_MICROSCOPICA,
    DIAGNOSTICO_HISTOPATOLOGICO
)

# --- Tests para matcheaCategoria ---

def test_matcheaCategoria_lineaAntecedentes_retornaTuplaActualizada():
    linea = "Material remitido - Antecedentes:"
    matchea, seccion, nueva_desc, nombre_tabla = matcheaCategoria(linea, None, False, "")
    
    assert matchea is True
    assert seccion == MATERIAL_REMITIDO_ANTECEDENTES
    assert nueva_desc is False

def test_matcheaCategoria_lineaMicroscopica_activaFlagNuevaDesc():
    linea = "Descripción microscópica"
    # Le pasamos "OtraSeccion" para simular que veníamos leyendo otra cosa
    matchea, seccion, nueva_desc, nombre_tabla = matcheaCategoria(linea, "OtraSeccion", False, "")
    
    assert matchea is True
    assert seccion == DESCRIPCION_MICROSCOPICA
    # Acá validamos que tu lógica levante la flag correctamente
    assert nueva_desc is True

def test_matcheaCategoria_lineaSinCategoria_mantieneEstadoActual():
    linea = "Este es un párrafo de texto normal describiendo la biopsia."
    estado_previo = DESCRIPCION_MACROSCOPICA
    
    matchea, seccion, nueva_desc, nombre_tabla = matcheaCategoria(linea, estado_previo, False, "TablaX")
    
    assert matchea is False
    assert seccion == estado_previo
    assert nombre_tabla == "TablaX"

# --- Tests para cargarDescMicro ---

def test_cargarDescMicro_bloqueValido_agregaAListaYReiniciaEstructura():
    lista_descripciones = []
    bloque_actual = {
        "Descripcion": "Células neoplásicas",
        "Diagnostico": {
            "Descripcion": "Carcinoma",
            "Imagenes": ["img1.png"]
        },
        "Tabla de Grado": []
    }
    
    bloque_reiniciado = cargarDescMicro(lista_descripciones, bloque_actual)
    
    # 1. Validamos que se haya agregado a la lista general
    assert len(lista_descripciones) == 1
    assert lista_descripciones[0]["Descripcion"] == "Células neoplásicas"
    
    # 2. Validamos que retorne un bloque en blanco para seguir procesando
    assert bloque_reiniciado["Descripcion"] == ""
    assert bloque_reiniciado["Diagnostico"]["Descripcion"] == ""
    assert len(bloque_reiniciado["Diagnostico"]["Imagenes"]) == 0

def test_cargarDescMicro_diagnosticoPuntoOVacio_loConvierteANone():
    lista_descripciones = []
    bloque_actual = {
        "Descripcion": "Tejido adiposo",
        "Diagnostico": {
            "Descripcion": ".", # Tu código tiene una regla para limpiar este punto
            "Imagenes": []
        },
        "Tabla de Grado": []
    }
    
    cargarDescMicro(lista_descripciones, bloque_actual)
    
    # Validamos que ese punto solitario se haya transformado en un None lógico
    assert lista_descripciones[0]["Diagnostico"]["Descripcion"] is None

# --- Tests para resultados ---

def test_resultados_seccionesConDatos_joineaStringsYCombinaDicts():
    datos_paciente = {"Protocolo": "CAN-001", "Especie": "Canino"}
    secciones = {
        MATERIAL_REMITIDO_ANTECEDENTES: ["Piel.", "Muestra de 3cm."],
        DESCRIPCION_MICROSCOPICA: [{"Descripcion": "Bloque microscópico"}]
    }
    
    resultado_final = resultados(datos_paciente, secciones)
    
    # Validamos la herencia de datos
    assert resultado_final["Protocolo"] == "CAN-001"
    
    # Validamos que los strings se unan con espacios
    assert resultado_final[MATERIAL_REMITIDO_ANTECEDENTES] == "Piel. Muestra de 3cm."
    
    # Validamos que la lista microscópica quede intacta (no se joinea)
    assert isinstance(resultado_final[DESCRIPCION_MICROSCOPICA], list)
    assert len(resultado_final[DESCRIPCION_MICROSCOPICA]) == 1

def test_resultados_seccionesVacias_asignaNoEncontrado():
    datos_paciente = {}
    secciones = {
        DESCRIPCION_MACROSCOPICA: [], # Lista vacía simulando que no encontró texto
        DESCRIPCION_MICROSCOPICA: []
    }
    
    resultado_final = resultados(datos_paciente, secciones)
    
    # Validamos que tu sistema de fallback "No encontrado" funcione
    assert resultado_final[DESCRIPCION_MACROSCOPICA] == "No encontrado"