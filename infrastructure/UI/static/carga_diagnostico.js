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
fileInput.addEventListener('change', function() {
      if (this.files && this.files[0]) processFile(this.files[0]);
});

function handleDrop(e) {
      let dt = e.dataTransfer;
      let files = dt.files;
      if(files.length > 0) {
            processFile(files[0]);
      }
}

function processFile(file) {
      nombre_archivo = file.name;
      dropText.style.display = 'none';
      dropIcon.style.display = 'none';
      dropSubtext.style.display = 'none';
      loader.style.display = 'block';

      // Animación simulada de IA extrayendo datos
      setTimeout(() => {
            loader.style.display = 'none';
            
            // Mostrar estado de éxito en el dropzone (en lugar de ocultarlo todo)
            dropIcon.innerHTML = '✅';
            dropIcon.style.display = 'block';
            dropText.innerHTML = '¡Documento procesado!';
            dropText.style.display = 'block';
            dropSubtext.innerHTML = file.name;
            dropSubtext.style.display = 'block';
            dropZone.style.borderColor = 'var(--success)';
            dropZone.style.backgroundColor = 'var(--success-bg)';

            // populateForm(mockExtractedData);
            
            // Si la pantalla es pequeña, hacer scroll automático al formulario
            if(window.innerWidth <= 900) {
                  document.getElementById('form-section').scrollIntoView({ behavior: 'smooth', block: 'start' });
            }
      }, 1500);
}

function populateForm(data) {
      document.getElementById('f-protocolo').value = data.Protocolo || '';
      document.getElementById('f-fecha').value = data.Fecha || '';
      document.getElementById('f-paciente').value = data.Paciente || '';
      document.getElementById('f-propietario').value = data.Propietario || '';
      document.getElementById('f-especie').value = data.Especie || '';
      document.getElementById('f-raza').value = data.Raza || '';
      document.getElementById('f-edad').value = data.Edad || '';
      document.getElementById('f-solicitante').value = data.Solicitante || '';
      document.getElementById('f-tecnica').value = data.Técnica || '';
      document.getElementById('f-mastocitomas').checked = data["Referencias mastocitomas"] === true;
      document.getElementById('f-antecedentes').value = data["Material remitido - Antecedentes"] || '';
      document.getElementById('f-macroscopica').value = data["Descripción macroscópica"] || '';

      const microContainer = document.getElementById('micro-container');
      microContainer.innerHTML = ''; 

      data["Descripción microscópica"].forEach((item, index) => {
            addMicroCard(item.Descripcion, item.Diagnostico.Descripcion, index);
      });
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
                  const confirmBox = confirm("Los datos han sido validados y guardados exitosamente.");
                  resetPage();
                  btn.innerHTML = 'Confirmar y Subir al Servidor';
                  btn.classList.add('btn-primary');
                  btn.classList.remove('btn-success');
                  btn.style.pointerEvents = 'auto';
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
    const fieldIds = ['f-protocolo','f-fecha','f-paciente','f-propietario','f-especie',
                      'f-raza','f-edad','f-solicitante','f-tecnica','f-antecedentes','f-macroscopica'];
    fieldIds.forEach(id => state.fields[id] = document.getElementById(id)?.value ?? '');
    state.fields['f-mastocitomas'] = document.getElementById('f-mastocitomas')?.checked ?? false;
    
    // Micro-cards: descripción, diagnóstico e imágenes por cada una
    document.querySelectorAll('.micro-card').forEach(card => {
        const desc = card.querySelector('textarea')?.value ?? '';
        const diag = card.querySelector('input[type="text"]')?.value ?? '';
        const images = Array.from(card.querySelectorAll('.image-thumbnail'))
                            .map(img => img.getAttribute('src'));
        state.microCards.push({ descripcion: desc, diagnostico: diag, imagenes: images });
    });
    
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
            alert('Se necesita información de diagnóstico para enviar al servidor');
      }
}
