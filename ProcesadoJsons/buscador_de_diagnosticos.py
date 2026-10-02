import json
import os
import re


# Patrones para extraer el grado y el diagnóstico
PATRON_GRADO = re.compile(
    r"\b(?:de\s+)?grado\s+(?P<grado>[IVXLCDM]+|\d+)\b",
    re.IGNORECASE
)

PATRON_DIAGNOSTICO = re.compile(
    r"\bcompatible\s+con\s+"
    r"(?P<diagnostico>.*?)(?="
    r"\s+(?:de\s+)?grado\s+[IVXLCDM]+"
    r"|\s+(?:de\s+)?grado\s+\d+"
    r"|\.|$)",
    re.IGNORECASE
)

def extraer_grado(texto):
    coincidencia = PATRON_GRADO.search(texto)

    if not coincidencia:
        return None

    grado = coincidencia.group("grado").upper()

    grados_romanos = {
        "I": 1,
        "II": 2,
        "III": 3,
        "IV": 4,
        "V": 5
    }

    if grado in grados_romanos:
        return grados_romanos[grado]

    return int(grado)


def extraer_diagnostico(texto):
    """
    Extrae el texto que aparece después de 'compatible con'.

    El diagnóstico termina:
    - antes de 'grado X', si existe; o
    - en el primer punto; o
    - al finalizar el texto.
    """
    coincidencia = PATRON_DIAGNOSTICO.search(texto)

    if not coincidencia:
        return None

    diagnostico = coincidencia.group("diagnostico")
    diagnostico = re.sub(r"\s+", " ", diagnostico)
    return diagnostico.strip(" ,;:-")


def extraer_datos_estructurados_de_json(ruta_archivo, nombre_archivo):
    """
    Extrae un resultado por cada descripción microscópica del JSON.

    Solo se leen los diagnósticos de:
        Descripción microscópica -> Diagnostico -> Descripcion

    Cada resultado contiene únicamente:
        diagnostico
        grado
    """
    with open(ruta_archivo, "r", encoding="utf-8") as archivo:
        data = json.load(archivo)

    resultados = []
    descripciones = data.get("Descripción microscópica", [])

    if not isinstance(descripciones, list):
        return resultados

    for bloque in descripciones:
        if not isinstance(bloque, dict):
            continue

        diagnostico = bloque.get("Diagnostico", {})
        if not isinstance(diagnostico, dict):
            continue

        texto = diagnostico.get("Descripcion", "")
        if not isinstance(texto, str) or not texto.strip():
            continue

        diagnostico_extraido = extraer_diagnostico(texto)
        if diagnostico_extraido is None:
            continue

        resultados.append({
            "diagnostico": diagnostico_extraido,
            "grado": extraer_grado(texto),
            "protocolo": data.get("Protocolo", "")
        })

    return resultados


if __name__ == "__main__":
    carpeta_jsons = "JSONS/"

    if not os.path.exists(carpeta_jsons):
        print(f"La carpeta {carpeta_jsons} no existe.")
        raise SystemExit(1)

    archivos_json = [
        archivo for archivo in os.listdir(carpeta_jsons)
        if archivo.lower().endswith(".json")
    ]

    if not archivos_json:
        print("No se encontraron archivos JSON en la carpeta.")
        raise SystemExit(0)

    resultados = []

    for nombre_archivo in archivos_json:
        ruta_archivo = os.path.join(carpeta_jsons, nombre_archivo)
        resultado = extraer_datos_estructurados_de_json(
            ruta_archivo,
            nombre_archivo
        )
        resultados.extend(resultado)
    with open("ProcesadoJsons/diagnosticos_extraidos.json", "w", encoding="utf-8") as archivo:
        json.dump(resultados, archivo, ensure_ascii=False, indent=2)

    print(resultados)