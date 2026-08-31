/*
 * Visor de imagenes a pantalla completa para las microfotografias.
 *
 * Los listeners se delegan en `document` porque el detalle del paciente
 * llega por HTMX y reemplaza #main-section: bindear por imagen al cargar
 * la pagina dejaria de funcionar tras el primer swap.
 */
(function () {
      "use strict";

      var overlay, imgEl, counterEl, btnPrev, btnNext, btnClose;
      var imagenes = [];     // srcs del grupo abierto
      var textos = [];       // alts del grupo abierto
      var indice = 0;
      var origen = null;     // .micro-image desde donde se abrio, para devolver el foco

      var MARKUP =
            '<button class="lightbox-close" aria-label="Cerrar" type="button">&times;</button>' +
            '<button class="lightbox-prev" aria-label="Anterior" type="button">&#10094;</button>' +
            '<img class="lightbox-img" src="" alt="">' +
            '<button class="lightbox-next" aria-label="Siguiente" type="button">&#10095;</button>' +
            '<div class="lightbox-counter"></div>';

      function cachearNodos() {
            // Se busca SIEMPRE en el documento actual: guardar la referencia
            // para siempre fallaba tras un swap de HTMX, y hay shells (la home,
            // servida desde index.html) que no traen el overlay en su markup.
            var actual = document.getElementById("lightbox");

            if (!actual) {
                  // Lo creamos al vuelo y lo colgamos del <body>, fuera de
                  // #main-section, para que sobreviva a los swaps.
                  actual = document.createElement("div");
                  actual.id = "lightbox";
                  actual.className = "lightbox";
                  actual.hidden = true;
                  actual.setAttribute("aria-hidden", "true");
                  actual.setAttribute("role", "dialog");
                  actual.setAttribute("aria-modal", "true");
                  actual.setAttribute("aria-label", "Visor de imagen");
                  actual.innerHTML = MARKUP;
                  document.body.appendChild(actual);
            }

            if (overlay === actual && imgEl) return overlay;

            overlay = actual;
            imgEl = overlay.querySelector(".lightbox-img");
            counterEl = overlay.querySelector(".lightbox-counter");
            btnPrev = overlay.querySelector(".lightbox-prev");
            btnNext = overlay.querySelector(".lightbox-next");
            btnClose = overlay.querySelector(".lightbox-close");
            return overlay;
      }

      // El overlay puede haber sido reemplazado por un swap de HTMX,
      // asi que preguntamos siempre por el nodo vivo.
      function abierto() {
            var n = document.getElementById("lightbox");
            return n && !n.hidden ? n : null;
      }

      function mostrar() {
            imgEl.src = imagenes[indice];
            imgEl.alt = textos[indice] || "";
            counterEl.textContent = (indice + 1) + " / " + imagenes.length;
            // con una sola imagen las flechas no aportan nada
            var varias = imagenes.length > 1;
            btnPrev.hidden = !varias;
            btnNext.hidden = !varias;
            counterEl.hidden = !varias;
      }

      function abrir(contenedor, clickeada) {
            if (!cachearNodos()) return;

            var nodos = contenedor.querySelectorAll(".micro-image img");
            if (!nodos.length) return;

            imagenes = [];
            textos = [];
            for (var i = 0; i < nodos.length; i++) {
                  imagenes.push(nodos[i].src);
                  textos.push(nodos[i].alt);
            }

            var propia = clickeada.querySelector("img");
            indice = 0;
            for (var j = 0; j < nodos.length; j++) {
                  if (nodos[j] === propia) { indice = j; break; }
            }

            origen = clickeada;
            mostrar();
            overlay.hidden = false;
            overlay.setAttribute("aria-hidden", "false");
            document.body.style.overflow = "hidden";
            btnClose.focus();
      }

      function cerrar() {
            var n = abierto();
            if (!n) {
                  // Aun cerrado, nunca dejamos el scroll bloqueado.
                  document.body.style.overflow = "";
                  return;
            }
            cachearNodos();
            n.hidden = true;
            n.setAttribute("aria-hidden", "true");
            if (imgEl) imgEl.src = "";
            document.body.style.overflow = "";
            if (origen && document.contains(origen)) origen.focus();
            origen = null;
      }

      function navegar(paso) {
            if (imagenes.length < 2) return;
            // wrap-around en ambos sentidos
            indice = (indice + paso + imagenes.length) % imagenes.length;
            mostrar();
      }

      document.addEventListener("click", function (e) {
            if (abierto()) {
                  if (e.target.closest(".lightbox-close")) { cerrar(); return; }
                  if (e.target.closest(".lightbox-prev")) { navegar(-1); return; }
                  if (e.target.closest(".lightbox-next")) { navegar(1); return; }
                  if (e.target === overlay) { cerrar(); return; }   // click en el backdrop
                  return;
            }

            var celda = e.target.closest(".micro-image");
            if (!celda) return;
            var grupo = celda.closest(".micro-images");
            if (!grupo) return;
            abrir(grupo, celda);
      });

      // Abrir con teclado desde la celda enfocada (role="button" + tabindex)
      document.addEventListener("keydown", function (e) {
            if (!abierto()) {
                  if (e.key !== "Enter" && e.key !== " ") return;
                  var celda = document.activeElement && document.activeElement.closest
                        ? document.activeElement.closest(".micro-image")
                        : null;
                  if (!celda) return;
                  var grupo = celda.closest(".micro-images");
                  if (!grupo) return;
                  e.preventDefault();
                  abrir(grupo, celda);
                  return;
            }

            if (e.key === "Escape") { cerrar(); }
            else if (e.key === "ArrowLeft") { navegar(-1); }
            else if (e.key === "ArrowRight") { navegar(1); }
      });

      // Si HTMX reemplaza el contenido con el visor abierto, lo cerramos:
      // si no, el overlay queda flotando sobre la pantalla nueva y el body
      // se queda con overflow:hidden (pagina sin scroll).
      // (el script corre en el <head>, asi que document.body todavia no existe)
      document.addEventListener("htmx:beforeSwap", cerrar);
})();
