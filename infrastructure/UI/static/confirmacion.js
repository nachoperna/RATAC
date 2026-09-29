// Reemplaza el confirm() del navegador por un dialogo con el estilo del sistema.
//
// htmx dispara htmx:confirm antes de cada pedido con hx-confirm. Al cancelar el
// evento y llamar despues a evt.detail.issueRequest() decidimos nosotros cuando
// continuar, sin tocar los hx-confirm que ya estan en las plantillas.

function textoDialogo(mensaje) {
      // Los hx-confirm vienen como "Titulo? Detalle." — se parte en la primera
      // pregunta para poder mostrar titulo y cuerpo por separado.
      const corte = mensaje.indexOf('? ');
      if (corte === -1) return { titulo: mensaje, cuerpo: '' };
      return {
            titulo: mensaje.slice(0, corte + 1).trim(),
            cuerpo: mensaje.slice(corte + 2).trim(),
      };
}

function mostrarConfirmacion(mensaje, alConfirmar) {
      const { titulo, cuerpo } = textoDialogo(mensaje);
      // Una accion irreversible se marca en rojo
      const esPeligro = /no se puede deshacer|eliminar|borrar/i.test(mensaje);

      const overlay = document.createElement('div');
      overlay.className = 'confirm-overlay';
      overlay.innerHTML = `
            <div class="confirm-content" role="alertdialog" aria-modal="true">
                  <div class="confirm-icon ${esPeligro ? 'confirm-icon-peligro' : ''}">
                        <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                              <circle cx="12" cy="12" r="10"></circle>
                              <line x1="12" y1="8" x2="12" y2="12"></line>
                              <line x1="12" y1="16" x2="12.01" y2="16"></line>
                        </svg>
                  </div>
                  <h3 class="confirm-title"></h3>
                  <p class="confirm-message"></p>
                  <div class="confirm-actions">
                        <button type="button" class="confirm-cancelar">Cancelar</button>
                        <button type="button" class="confirm-aceptar ${esPeligro ? 'confirm-aceptar-peligro' : ''}">${esPeligro ? 'Eliminar' : 'Aceptar'}</button>
                  </div>
            </div>`;

      // Como texto y no como HTML: el mensaje puede incluir datos cargados por el usuario
      overlay.querySelector('.confirm-title').textContent = titulo;
      const parrafo = overlay.querySelector('.confirm-message');
      if (cuerpo) {
            parrafo.textContent = cuerpo;
      } else {
            parrafo.remove();
      }

      document.body.appendChild(overlay);
      requestAnimationFrame(() => overlay.classList.add('active'));

      const cerrar = () => {
            overlay.classList.remove('active');
            document.removeEventListener('keydown', alPresionarTecla);
            setTimeout(() => overlay.remove(), 250);
      };

      function alPresionarTecla(e) {
            if (e.key === 'Escape') cerrar();
      }

      overlay.querySelector('.confirm-cancelar').addEventListener('click', cerrar);
      overlay.querySelector('.confirm-aceptar').addEventListener('click', () => {
            cerrar();
            alConfirmar();
      });
      // Clic fuera del recuadro = cancelar
      overlay.addEventListener('click', (e) => { if (e.target === overlay) cerrar(); });
      document.addEventListener('keydown', alPresionarTecla);

      overlay.querySelector('.confirm-aceptar').focus();
}

document.addEventListener('htmx:confirm', (evt) => {
      if (!evt.detail.question) return; // el pedido no pide confirmacion
      evt.preventDefault();             // frenamos el confirm() nativo
      mostrarConfirmacion(evt.detail.question, () => evt.detail.issueRequest(true));
});

// Aviso de una accion que ya ocurrio (por ejemplo, un guardado exitoso).
// A diferencia de la confirmacion no ofrece "Cancelar": no hay nada que
// cancelar, el cambio ya esta hecho.
function mostrarAviso(mensaje, titulo, alCerrar) {
      const overlay = document.createElement('div');
      overlay.className = 'confirm-overlay';
      overlay.innerHTML = `
            <div class="confirm-content" role="alertdialog" aria-modal="true">
                  <div class="confirm-icon confirm-icon-exito">
                        <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                              <polyline points="22 4 12 14.01 9 11.01"></polyline>
                        </svg>
                  </div>
                  <h3 class="confirm-title"></h3>
                  <p class="confirm-message"></p>
                  <div class="confirm-actions confirm-actions-unico">
                        <button type="button" class="confirm-aceptar">Aceptar</button>
                  </div>
            </div>`;

      overlay.querySelector('.confirm-title').textContent = titulo || 'Listo';
      overlay.querySelector('.confirm-message').textContent = mensaje;

      document.body.appendChild(overlay);
      requestAnimationFrame(() => overlay.classList.add('active'));

      const cerrar = () => {
            overlay.classList.remove('active');
            document.removeEventListener('keydown', alPresionarTecla);
            setTimeout(() => overlay.remove(), 250);
            if (alCerrar) alCerrar();
      };

      function alPresionarTecla(e) {
            if (e.key === 'Escape' || e.key === 'Enter') cerrar();
      }

      overlay.querySelector('.confirm-aceptar').addEventListener('click', cerrar);
      overlay.addEventListener('click', (e) => { if (e.target === overlay) cerrar(); });
      document.addEventListener('keydown', alPresionarTecla);

      overlay.querySelector('.confirm-aceptar').focus();
}
