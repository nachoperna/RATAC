package ui

import (
	"RATAC/application"
	"RATAC/views"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Mensajes que se le muestran al usuario. Son deliberadamente generales: el
// detalle tecnico del error va al log del servidor, no a la pantalla, porque
// el texto crudo de un error expone rutas, nombres de tablas y datos internos.
const (
	MSJ_ERROR_GUARDAR   = "No se pudo guardar el diagnóstico. No se registró ningún cambio, puede intentarlo nuevamente."
	MSJ_ERROR_DUPLICADO = "El diagnóstico subido ya se encuentra cargado en el sistema."
	MSJ_ERROR_ARCHIVO   = "No se pudo procesar el archivo. Verifique que sea un PDF o DOCX válido."
	MSJ_ERROR_GENERICO  = "No se pudo completar la operación. Inténtelo nuevamente."
	MSJ_ERROR_SESION    = "Se venció su sesión en RATAC. Deberá ingresar nuevamente."
	MSJ_ERROR_ELIMINAR  = "No se pudo eliminar el diagnóstico. Inténtelo nuevamente."
	MSJ_ERROR_PROTOCOLO = "El documento no incluye el número de protocolo. Complételo en el formulario antes de confirmar la carga."

	MSJ_ERROR_REGISTRO        = "No se pudo enviar la solicitud. Inténtelo nuevamente."
	// El mismo texto para un correo ya registrado y para cualquier otro fallo del
	// alta: distinguirlos permitiria averiguar que correos tienen cuenta probando
	// direcciones en el formulario.
	MSJ_ERROR_REGISTRO_EMAIL  = "No se pudo completar el registro. Verifique los datos ingresados o ingrese desde la pestaña «Ingresar» si ya tiene una cuenta."
	MSJ_ERROR_REGISTRO_CAMPOS = "Complete el nombre del laboratorio, el correo y la ciudad antes de enviar la solicitud."

	MSJ_ERROR_BUSQUEDA     = "No se pudo completar la búsqueda. Inténtelo nuevamente."
	MSJ_SIN_MAS_RESULTADOS = "No hay más resultados para mostrar."
)

// errorAlUsuario responde con un mensaje entendible y deja el detalle real en
// el log del servidor.
func errorAlUsuario(w http.ResponseWriter, status int, mensaje string, detalle error, contexto string) {
	if detalle != nil {
		log.Printf("[%s] %v", contexto, detalle)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(mensaje))
}

// mensajeDocumento traduce el fallo al procesar un documento en el titulo y el
// texto que se le muestra al usuario. Cada causa tiene su propio mensaje: antes
// todas se informaban como "ya se encuentra cargado", que solo era cierto para
// el duplicado.
func mensajeDocumento(err error, nombreArchivo string) (titulo, mensaje string) {
	switch {
	case errors.Is(err, application.ErrDiagnosticoRepetido):
		return "Diagnóstico duplicado", "El diagnóstico subido ya se encuentra cargado en el sistema."
	case errors.Is(err, application.ErrTipoArchivoInvalido):
		return "Formato no admitido", fmt.Sprintf("El archivo «%s» no es un PDF ni un DOCX. Suba el documento en alguno de esos formatos.", nombreArchivo)
	case errors.Is(err, application.ErrDocumentoIlegible):
		return "No se pudo leer el archivo", fmt.Sprintf("No se pudo leer el contenido de «%s». Verifique que el archivo no esté dañado.", nombreArchivo)
	case errors.Is(err, application.ErrFormatoDesconocido):
		return "No se pudo interpretar el documento", fmt.Sprintf("«%s» no tiene la estructura de un diagnóstico. Revise que sea el documento correcto y que conserve sus secciones y tablas.", nombreArchivo)
	default:
		return "", MSJ_ERROR_ARCHIVO
	}
}

// mensajeInsercion elige que avisar cuando falla el alta de un diagnostico.
// Un protocolo faltante lo puede resolver el usuario completando el campo; el
// resto de los fallos no, y se informan de forma generica.
func mensajeInsercion(err error) string {
	switch {
	case errors.Is(err, application.ErrProtocoloFaltante):
		return MSJ_ERROR_PROTOCOLO
	case errors.Is(err, application.ErrProtocoloDuplicado):
		return "Ya existe un diagnóstico cargado con ese número de protocolo. Verifique el número ingresado."
	default:
		return MSJ_ERROR_GUARDAR
	}
}

// errorInline muestra el aviso dentro de la pantalla, sin interrumpir al usuario.
// Se responde 200 porque htmx descarta el cuerpo de las respuestas de error y el
// cartel nunca llegaria a insertarse.
func errorInline(w http.ResponseWriter, r *http.Request, mensaje string, detalle error, contexto string) {
	if detalle != nil {
		log.Printf("[%s] %v", contexto, detalle)
	}
	// Si el pedido no viene de htmx (por ejemplo una recarga con F5 sobre la URL de
	// una busqueda) devolver solo el fragmento del aviso mostraria una pagina rota:
	// se vuelve al inicio.
	if r.Header.Get("HX-Request") != "true" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// El pedido que falla suele apuntar a un hx-target con outerHTML (por ejemplo
	// #main-section). Se redirige el swap a un destino inofensivo para que la
	// seccion no se reemplace con el fragmento del aviso, que entra por su
	// hx-swap-oob. Con HX-Reswap "none" htmx 1.9 descarta la respuesta entera y el
	// aviso tampoco se aplicaria.
	// El aviso se inserta directamente en su contenedor. Antes se redirigia el
	// swap a un div oculto y el aviso entraba por hx-swap-oob, pero htmx no
	// procesa el oob cuando el destino del swap esta oculto.
	w.Header().Set("HX-Retarget", "#error-inline")
	w.Header().Set("HX-Reswap", "outerHTML")
	// Sin esto, el hx-push-url del buscador deja en la barra de direcciones la URL
	// del pedido que fallo; al recargar con F5 el navegador pide esa ruta y recibe
	// el fragmento del aviso suelto, sin la pagina.
	w.Header().Set("HX-Push-Url", "false")
	// Se responde 422 y no 200: htmx solo aplica HX-Retarget en respuestas de
	// error, y sin el swap iria al destino original (la grilla de resultados).
	// El fragmento igual se muestra porque index.js habilita el swap de las
	// respuestas de error en htmx:beforeSwap.
	w.WriteHeader(http.StatusUnprocessableEntity)
	views.ErrorInline(mensaje).Render(r.Context(), w)
}
