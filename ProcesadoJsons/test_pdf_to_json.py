import pytest
from PDF_to_json import limpiar_celda, limpiar_diagnostico, procesar_raza_edad, reordenar_datos

# --- Tests para limpiar_celda ---

def test_limpiar_celda_conSaltosDeLinea_retornaStringLimpio():
    celda = "Mestizo\n(Canino)"
    resultado = limpiar_celda(celda)
    assert resultado == "Mestizo (Canino)"

def test_limpiar_celda_vaciaONula_retornaStringVacio():
    assert limpiar_celda(None) == ""
    assert limpiar_celda("   ") == ""

# --- Tests para limpiar_diagnostico ---

def test_limpiar_diagnostico_conFirmasVeterinarias_retornaTextoLimpio():
    # Simulamos cómo a veces el PDF pega la firma al final del diagnóstico
    desc_sucia = "Mastocitoma grado II. Laura Denzoin Médica Veterinaria M.P. 1234"
    resultado = limpiar_diagnostico(desc_sucia)
    
    assert resultado == "Mastocitoma grado II."
    assert "Laura Denzoin" not in resultado

def test_limpiar_diagnostico_conMultiplesEspacios_retornaTextoConEspaciosSimples():
    desc_sucia = "Carcinoma    de células   escamosas."
    resultado = limpiar_diagnostico(desc_sucia)
    
    assert resultado == "Carcinoma de células escamosas."

# --- Tests para procesar_raza_edad ---

def test_procesar_raza_edad_stringConAmbosDatos_retornaDiccionarioActualizado():
    # Simulamos el string que llega de la cabecera
    valor_entrada = "Caniche - 10 años"
    datos_iniciales = {}
    
    resultado = procesar_raza_edad(valor_entrada, datos_iniciales)
    
    assert resultado["Raza"] == "Caniche"
    assert resultado["Edad"] == "10"

def test_procesar_raza_edad_stringSoloRaza_retornaDiccionarioSinEdad():
    valor_entrada = "Felino Europeo"
    datos_iniciales = {}
    
    resultado = procesar_raza_edad(valor_entrada, datos_iniciales)
    
    assert resultado["Raza"] == "Felino Europeo"
    assert "Edad" not in resultado  # Como no hay números, no debería crear la key Edad

# --- Tests para reordenar_datos ---

def test_reordenar_datos_datosDesordenadosYVacios_retornaDiccionarioOrdenadoYNulos():
    datos_desordenados = {
        "Especie": "Felino",
        "Protocolo": "FEL-2026",
        "Raza": "-",      # Debería pasar a None (null en json)
        "Paciente": "",   # Debería pasar a None
        "DatoExtra": "ignorado por el orden pero conservado al final"
    }
    
    resultado = reordenar_datos(datos_desordenados)
    
    # Verificamos que el orden sea el que exige el frontend/BD
    claves = list(resultado.keys())
    assert claves[0] == "Protocolo"
    assert claves[5] == "Especie"
    
    # Verificamos la limpieza de vacíos a None
    assert resultado["Raza"] is None
    assert resultado["Paciente"] is None
    assert resultado["DatoExtra"] == "ignorado por el orden pero conservado al final"