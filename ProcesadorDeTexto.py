import os
import json
import re
import spacy
from sklearn.feature_extraction.text import CountVectorizer

# 1. Cargar el modelo de lenguaje en español
print("Cargando modelo de spaCy (es_core_news_sm)...")
nlp = spacy.load("es_core_news_sm")

def extraer_datos_estructurados_de_json(ruta_archivo, nombre_archivo):
    """
    Lee el JSON y devuelve una lista de diccionarios, separando cada diagnóstico
    junto con su ID de documento y el grado tumoral extraído.
    """
    with open(ruta_archivo, 'r', encoding='utf-8') as f:
        data = json.load(f)


    id_doc = data.get("Protocolo", nombre_archivo)
    resultados = []

    micro = data.get("Descripción microscópica", [])
    if isinstance(micro, list):
        for bloque in micro:
            diagnostico = bloque.get("Diagnostico", {})
            if isinstance(diagnostico, dict):
                desc_diag = diagnostico.get("Descripcion", "")
                if desc_diag:
                    # Expresión Regular para buscar el Grado
                    match = re.search(r"(?i)\bgrado\s*([IVX]+|\d)\b", desc_diag)
                    grado_encontrado = match.group(1).upper() if match else "No especificado"

                    resultados.append({
                        "id_documento": id_doc,
                        "texto_crudo": desc_diag,
                        "grado": grado_encontrado
                    })

    return resultados

def preprocesar_texto(texto):
    """
    PASO 2: Filtrado Gramatical (POS Tagging).
    Limpia el texto y acepta ÚNICAMENTE sustantivos y adjetivos.
    """
    if not texto.strip():
        return ""
        
    doc = nlp(texto)
    tokens_lematizados = []
    
    for token in doc:
        if not token.is_stop and not token.is_punct and not token.is_space and token.is_alpha:
            # Solo dejamos pasar Sustantivos y Adjetivos
            #if token.pos_ in ["NOUN", "ADJ"]:
            tokens_lematizados.append(token.lemma_.lower())
            
    return " ".join(tokens_lematizados)

def extraer_ngramas_por_frecuencia(corpus_procesado):
    """
    Extrae n-gramas del corpus completo.
    """
    vectorizador = CountVectorizer(
        ngram_range=(2, 4), # Extraemos n-gramas de 2 a 4 palabras
        # min_df=min_doc_freq,
        min_df= 1,
        max_df=0.8
    )
    
    try:
        vectorizador.fit(corpus_procesado)
        return vectorizador.get_feature_names_out()
    except ValueError:
        return ["No se encontraron n-gramas que cumplan los filtros."]

# --- Ejecución del Pipeline ---

if __name__ == "__main__":
    carpeta_jsons = "JSONS/"
    
    if not os.path.exists(carpeta_jsons):
        print(f"La carpeta {carpeta_jsons} no existe. Por favor corre primero diagToJson.py")
        exit()

    archivos_json = [f for f in os.listdir(carpeta_jsons) if f.endswith('.json')]
    
    if not archivos_json:
        print("No se encontraron archivos JSON en la carpeta.")
        exit()

    corpus_limpio = []
    metadatos_diagnosticos = [] # Datos para la base

    print(f"\n--- 1. Procesando {len(archivos_json)} documentos JSON ---")
    
    for archivo in archivos_json:
        ruta_completa = os.path.join(carpeta_jsons, archivo)
        
        # Extraemos la lista estructurada en lugar del texto crudo unificado
        diagnosticos_doc = extraer_datos_estructurados_de_json(ruta_completa, archivo)
        
        # B. Iteramos sobre cada diagnóstico individualmente
        for diag in diagnosticos_doc:
            texto_limpio = preprocesar_texto(diag["texto_crudo"])
            
            if texto_limpio:
                corpus_limpio.append(texto_limpio)
                
                # Guardamos la relación completa para la futura base de datos
                metadatos_diagnosticos.append({
                    "id_documento": diag["id_documento"],
                    "grado": diag["grado"],
                    "texto_procesado": texto_limpio
                })
                print(f"[OK] {diag['id_documento']} - Grado capturado: {diag['grado']}")

    print("\n--- 2. Extrayendo N-gramas del Corpus Completo ---")
    if corpus_limpio:
        ngramas_filtrados = extraer_ngramas_por_frecuencia(corpus_limpio)
        print(f"\nSe encontraron {len(ngramas_filtrados)} n-gramas relevantes:")
        for ngrama in ngramas_filtrados[:50]:
            print(f"- {ngrama}")