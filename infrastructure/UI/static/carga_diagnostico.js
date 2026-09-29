const dropZone = document.getElementById('drop-zone');
const fileInput = document.getElementById('file-input');
const dropText = document.getElementById('drop-zone-text');
const dropIcon = document.getElementById('drop-icon');
const dropSubtext = document.getElementById('drop-subtext');
const loader = document.getElementById('loader');
const imagenes_subidas = []
let nombre_archivo;
let estado_post_carga;
let estado_pre_envio;

['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
      dropZone.addEventListener(eventName, preventDefaults, false);
});

function preventDefaults(e) {
      e.preventDefault();
}

['dragenter', 'dragover'].forEach(eventName => {
      dropZone.addEventListener(eventName, () => dropZone.classList.add('dragover'), false);
});

['dragleave', 'drop'].forEach(eventName => {
      dropZone.addEventListener(eventName, () => dropZone.classList.remove('dragover'), false);
});

dropZone.addEventListener('drop', handleDrop, false);
dropZone.addEventListener('click', () => fileInput.click());

function handleDrop(e) {
      let dt = e.dataTransfer;
      let files = dt.files;
      if(files.length > 0) {
            // Asignar archivo al fileInput para que htmx lo detecte
            const dataTransfer = new DataTransfer();
            dataTransfer.items.add(files[0]);
            fileInput.files = dataTransfer.files;
            // Disparar change para que htmx envíe el POST
            fileInput.dispatchEvent(new Event('change', { bubbles: true }));
      }
}

function processFile(file) {
      nombre_archivo = file.name;
      dropText.style.display = 'none';
      dropIcon.style.display = 'none';
      dropSubtext.style.display = 'none';
      loader.style.display = 'block';

      setTimeout(() => {
            loader.style.display = 'none';
            
            dropIcon.innerHTML = '✅';
            dropIcon.style.display = 'block';
            dropText.innerHTML = '¡Documento procesado!';
            dropText.style.display = 'block';
            dropSubtext.innerHTML = file.name;
            dropSubtext.style.display = 'block';
            dropZone.style.borderColor = 'var(--success)';
            dropZone.style.backgroundColor = 'var(--success-bg)';
            
            // Si la pantalla es pequeña, hacer scroll automático al formulario
            if(window.innerWidth <= 900) {
                  document.getElementById('form-section').scrollIntoView({ behavior: 'smooth', block: 'start' });
            }
      }, 500);
}

function addMicroCard(descripcion = '', diagnostico = '', index = null) {
      const microContainer = document.getElementById('micro-container');
      const currentIndex = index !== null ? index : microContainer.children.length;
      const card = document.createElement('div');
      card.className = 'micro-card';

      card.innerHTML = `
            <div class="micro-card-header" style="display: flex; justify-content: space-between; align-items: center;">
                  <span>Muestra / Hallazgo #${currentIndex + 1}</span>
                  <span style="cursor: pointer; color: #d9534f; font-family: 'Source Serif 4', serif; font-size: 14px; text-transform: none; letter-spacing: 0; font-weight: normal;" onclick="this.closest('.micro-card').remove()">Eliminar</span>
            </div>
            <div class="form-group full-width" style="margin-bottom: 16px;">
                  <label class="form-label">Descripción Microscópica</label>
                  <textarea class="form-control">${descripcion}</textarea>
            </div>
            <div class="form-group full-width">
                  <label class="form-label">Diagnóstico</label>
                  <input type="text" class="form-control" value="${diagnostico}">
            </div>
            <div class="micro-images-section">
                  <div class="micro-images-header">
                        <span>Imágenes Asociadas</span>
                  </div>
                  <div class="image-preview-container">
                  <!-- Las miniaturas irán aquí -->
                  </div>
                  <button type="button" class="add-image-btn" onclick="this.nextElementSibling.click()">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
                        Añadir Imagen
                  </button>
                  <input type="file" style="display: none;" accept="image/*" multiple onchange="handleImageSelection(this)">
            </div>
      `;
      microContainer.appendChild(card);
}

// Manejador para previsualizar las imágenes seleccionadas
function handleImageSelection(input) {
      if (!input.files || input.files.length === 0) return;
      const container = input.previousElementSibling.previousElementSibling; // el div .image-preview-container

      Array.from(input.files).forEach(file => {
            const url = URL.createObjectURL(file);
            const entry = { file, url };
            imagenes_subidas.push(entry);


            const wrapper = document.createElement('div');
            wrapper.className = 'image-thumbnail-wrapper';
            wrapper.innerHTML = `
                  <img src="${url}" class="image-thumbnail" alt="Miniatura">
                  <button type="button" class="image-remove-btn" title="Eliminar imagen" onclick="this.parentElement.remove()">×</button>
            `;

            wrapper.querySelector('.image-remove-btn').addEventListener('click', () => {
                  URL.revokeObjectURL(url);
                  const idx = imagenes_subidas.indexOf(entry);
                  if (idx !== -1) imagenes_subidas.splice(idx, 1);
                  wrapper.remove();
            });

            container.appendChild(wrapper);
      });

      // Limpiamos el input para permitir seleccionar la misma imagen si se borra y se vuelve a añadir
      input.value = '';
}

function addEmptyMicroCard() {
      addMicroCard();
}

function resetPage() {
      // Resetear dropzone a estado inicial
      dropIcon.innerHTML = '📄';
      dropText.innerHTML = 'Arrastra tu documento aquí o haz clic para explorar';
      dropSubtext.innerHTML = 'Soporta .PDF, .DOCX';
      dropZone.style.borderColor = 'var(--lavender-mid)';
      dropZone.style.backgroundColor = '#fafbfc';
      fileInput.value = ''; 

      // Limpiar formulario
      document.querySelectorAll('#form-section input[type="text"], #form-section textarea').forEach(el => el.value = '');
      document.getElementById('f-mastocitomas').checked = false;
      document.getElementById('micro-container').innerHTML = '';
      addEmptyMicroCard(); 

      window.scrollTo({ top: 0, behavior: 'smooth' });
}

function submitData() {
      const btn = document.getElementById('btn-submit');
      btn.innerHTML = 'Subiendo...';
      btn.style.opacity = '0.8';
      btn.style.pointerEvents = 'none';

      setTimeout(() => {
            btn.innerHTML = '✓ ¡Datos Guardados!';
            btn.classList.add('btn-success');
            btn.classList.remove('btn-primary');
            btn.style.opacity = '1';
 
            setTimeout(() => {
                  // El formulario se limpia recien al cerrar el aviso, para que el
                  // usuario lo lea antes de que la pantalla cambie.
                  mostrarAviso(
                        'El diagnóstico fue validado y guardado correctamente.',
                        'Diagnóstico guardado',
                        () => {
                              resetPage();
                              btn.innerHTML = 'Confirmar y Subir al Servidor';
                              btn.classList.add('btn-primary');
                              btn.classList.remove('btn-success');
                              btn.style.pointerEvents = 'auto';
                        }
                  );
            }, 1000);
      }, 1500);
}

function getFileName(){ return nombre_archivo; }

function getNumeroDescMicro(){
      return document.getElementById('micro-container').children.length + 1;
}

function getRutaImagenes(){
      return Array.from(document.querySelectorAll('.image-thumbnail')).map(img => img.getAttribute('src'));
}

function capturarEstado() {
      const state = { fields: {}, microCards: [], images: [] };
      
      // Campos simples
      const fieldIds = ['f-protocolo','f-fecha','f-paciente','f-familia','f-especie',
                        'f-raza','f-edad','f-solicitante','f-tecnica','f-antecedentes','f-macroscopica'];
      fieldIds.forEach(id => state.fields[id] = document.getElementById(id)?.value ?? '');
      state.fields['f-mastocitomas'] = String(document.getElementById('f-mastocitomas')?.checked ?? false);
      
      // Micro-cards: descripción, diagnóstico e imágenes por cada una
      document.querySelectorAll('.micro-card').forEach(card => {
            const desc = card.querySelector('textarea')?.value ?? '';
            const diag = card.querySelector('input[type="text"]')?.value ?? '';
            const images = Array.from(card.querySelectorAll('.image-thumbnail'))
                              .map(img => img.getAttribute('src'));
            state.microCards.push({ descripcion: desc, diagnostico: diag, imagenes: images });
      });

      // Guardamos todas las imagenes juntas para comparar con las existentes en server
      state.images = getRutaImagenes();

      return state;
}

function guardarEstado(){
      estado_post_carga = capturarEstado();
}

function hayCambios() {
    if (!estado_post_carga) return true;
    estado_pre_envio = capturarEstado();
    return JSON.stringify(estado_post_carga) !== JSON.stringify(estado_pre_envio);
}

function hayDatos(){
      const archivo_subido = document.getElementById('file-input').files.length > 0;
      const protocolo_valor = document.getElementById('f-protocolo').value.trim() !== '';
      return archivo_subido || protocolo_valor;
}

function validacionDatosMinimos(e){
      if (!hayDatos()){
            e.preventDefault();
            mostrarError('Cargue un documento o complete al menos el número de protocolo antes de continuar.', 'Faltan datos del diagnóstico', 'aviso');
      }
}

function validacionYFormData(evt) {
      if (!hayDatos()) {
            evt.preventDefault();
            mostrarError('Cargue un documento o complete al menos el número de protocolo antes de continuar.', 'Faltan datos del diagnóstico', 'aviso');
            return;
      }
      // El protocolo identifica al diagnostico y es unico en la base. Hay documentos
      // que no lo traen adentro (queda vacio al extraer), asi que se pide completarlo
      // antes de enviar: sin el, la carga falla recien al insertar.
      const protocolo = document.getElementById('f-protocolo');
      if (protocolo && protocolo.value.trim() === '') {
            evt.preventDefault();
            mostrarError('El documento no incluye el número de protocolo. Complételo en el formulario antes de confirmar la carga.', 'Falta el número de protocolo', 'aviso');
            protocolo.focus();
            if (protocolo.scrollIntoView) {
                  protocolo.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
            return;
      }
      evt.preventDefault(); // cancela el request de htmx

      const estado = capturarEstado(); // siempre obtener fresco

      const reemplazarBlob = (src) => {
            const entry = imagenes_subidas.find(e => e.url === src);
            return entry ? entry.file.name : src;
      };

      estado.images = estado.images.map(reemplazarBlob);
      estado.microCards.forEach(card => {
            card.imagenes = card.imagenes.map(reemplazarBlob);
      });

      const formData = new FormData();
      formData.append('archivo', nombre_archivo);
      formData.append('cambios', String(hayCambios()));
      formData.append('campos', JSON.stringify(estado)); // usar estado, no capturarEstado()
      imagenes_subidas.forEach(e => formData.append('imagenes', e.file));

      fetch('/diagnosticos/alta/carga', { method: 'POST', body: formData })
            .then(async r => {
                  if (r.ok) {
                        submitData();
                        return;
                  }
                  if (r.status === 401) {
                        mostrarError('Se venció su sesión en RATAC. Deberá ingresar nuevamente.');
                        setTimeout(() => { window.location.href = '/ingreso/formulario'; }, 2500);
                        return;
                  }
                  // Cualquier otro fallo: antes no se avisaba nada y el usuario
                  // creia que el diagnostico se habia guardado.
                  const mensaje = await r.text().catch(() => '');
                  mostrarError(mensaje.trim() || 'No se pudo guardar el diagnóstico. No se registró ningún cambio, puede intentarlo nuevamente.');
            })
            .catch(err => {
                  // Falla de red: el request nunca llego al servidor
                  console.error('Error de red al subir el diagnóstico:', err);
                  mostrarError('No se pudo contactar al servidor. Verifique su conexión e inténtelo nuevamente.');
            });
}

function closeError() {
      document.getElementById('error').classList.remove('active');
}

// El aviso tambien se cierra con Escape o clickeando fuera del recuadro
document.addEventListener('keydown', (e) => {
      if (e.key !== 'Escape') return;
      const overlay = document.getElementById('error');
      if (overlay && overlay.classList.contains('active')) closeError();
});

document.addEventListener('click', (e) => {
      const overlay = document.getElementById('error');
      if (overlay && overlay.classList.contains('active') && e.target === overlay) closeError();
});

// Muestra un aviso reutilizando el overlay que ya usa el servidor
// (views/error_carga.templ). Sirve para los avisos que detecta el navegador,
// donde no hay una respuesta del servidor para insertar.
//
// tipo: 'error' (algo fallo) o 'aviso' (falta completar datos). Cambia solo el
// icono y su color; un dato faltante no deberia verse tan grave como un fallo.
function mostrarError(mensaje, titulo, tipo) {
      const overlay = document.getElementById('error');
      if (!overlay) {
            alert(mensaje); // ultimo recurso: la pantalla no tiene el contenedor
            return;
      }

      const esAviso = tipo === 'aviso';
      const icono = esAviso
            ? `<circle cx="12" cy="12" r="10"></circle>
               <line x1="12" y1="16" x2="12" y2="12"></line>
               <line x1="12" y1="8" x2="12.01" y2="8"></line>`
            : `<circle cx="12" cy="12" r="10"></circle>
               <line x1="12" y1="8" x2="12" y2="12"></line>
               <line x1="12" y1="16" x2="12.01" y2="16"></line>`;

      overlay.innerHTML = `
            <div class="error-content">
                  <div class="error-icon ${esAviso ? 'error-icon-aviso' : ''}">
                        <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                              ${icono}
                        </svg>
                  </div>
                  <h3 class="error-title" id="error-title"></h3>
                  <p class="error-message" id="error-message"></p>
                  <div class="error-actions">
                        <button class="btn btn-primary" onclick="closeError()">Aceptar</button>
                  </div>
            </div>`;
      // Se asigna como texto (no como HTML) para que el mensaje del servidor
      // no pueda inyectar marcado en la pagina.
      overlay.querySelector('#error-title').textContent = titulo || 'No se pudo completar la operación';
      overlay.querySelector('#error-message').textContent = mensaje;
      overlay.classList.add('active');

      // El foco va al boton para poder cerrar con Enter o Espacio
      const boton = overlay.querySelector('.error-actions .btn');
      if (boton) boton.focus();
}
