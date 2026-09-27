const searchInput = document.getElementById('citySearch');
const suggestionsList = document.getElementById('suggestionsList');
const submitBtn = document.getElementById('submitBtn');

document.addEventListener('DOMContentLoaded', () => {
    const tab = new URLSearchParams(window.location.search).get('tab');
    if (tab === 'collab') {
        const btn = document.getElementById('tab-btn-collab');
        if (btn) switchTab('collab', btn);  // quita active de "Ingresar" y lo pone en "Solicitar Colaboración"
    }
});

// Función para cambiar de pestaña
function switchTab(tabId, btnElement) {
      // 1. Quitar clase active de todos los botones y contenedores
      document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
      document.querySelectorAll('.form-container').forEach(form => form.classList.remove('active'));

      // 2. Añadir clase active al botón presionado y al formulario correspondiente
      btnElement.classList.add('active');
      document.getElementById('form-' + tabId).classList.add('active');
}

// Función para añadir una nueva fila de veterinario
function addVetRow() {
      const container = document.getElementById('vets-container');

      const row = document.createElement('div');
      row.className = 'vet-row';

      row.innerHTML = `
            <div class="vet-input-group">
                  <label class="form-label">Nombre y Apellido</label>
                  <input type="text" name="nombre-vet" class="form-control" placeholder="Dr/Dra. Nombre" required>
            </div>
            <div class="vet-input-group">
                  <label class="form-label">Matrícula (MP/MN)</label>
                  <input type="text" name="matricula-vet" class="form-control" placeholder="N° de Matrícula" required>
            </div>
            <button type="button" class="btn-remove-vet" title="Eliminar veterinario" onclick="removeVetRow(this)">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"></path></svg>
            </button>
      `;

      container.appendChild(row);
}

// Función para eliminar una fila de veterinario
function removeVetRow(btnElement) {
      const row = btnElement.closest('.vet-row');
      // Animación de salida antes de eliminar
      row.style.opacity = '0';
      row.style.transform = 'translateY(-10px)';
      row.style.transition = 'all 0.2s';

      setTimeout(() => {
            row.remove();
      }, 200);
}

function mismaContraseña(event) {
      const password = document.getElementById('new-pwd');
      const confirmPassword = document.getElementById('new-pwd-confirm');
      const errorMsg = document.getElementById('password-error');
      
      if (confirmPassword){
            // Comprobar si los valores son diferentes
            if (password.value !== confirmPassword.value) {
                  // Prevenir que el formulario se envíe
                  event.preventDefault(); 

                  // Mostrar el mensaje de error
                  errorMsg.style.display = 'block';

                  // (Opcional) Pintar los bordes de rojo para resaltar el error
                  password.style.borderColor = '#d32f2f';
                  confirmPassword.style.borderColor = '#d32f2f';
            } else {
                  // Si coinciden, ocultar el error y restaurar bordes (por si acaso)
                  errorMsg.style.display = 'none';
                  password.style.borderColor = '';
                  confirmPassword.style.borderColor = '';
            }
      }
}

// Función para alternar el modo de Cambiar Contraseña dentro del Login
function toggleChangePasswordMode(isChanging) {
      const labelPwd = document.getElementById('label-password');
      const linkChange = document.getElementById('link-change-pwd');
      const newFields = document.getElementById('new-password-fields');
      const btnSubmit = document.getElementById('btn-login-submit');
      const loginTitle = document.getElementById('login-title');
      const loginSubtitle = document.getElementById('login-subtitle');

      if (isChanging) {
            // Modo: Cambiar Contraseña
            labelPwd.textContent = 'Contraseña Actual';
            linkChange.style.display = 'none';
            
            // Mostrar campos extra con animación
            newFields.style.display = 'block';
            setTimeout(() => { newFields.style.opacity = '1'; }, 10);
            
            // Actualizar textos
            btnSubmit.textContent = 'Cambiar Contraseña';
            loginTitle.textContent = 'Cambiar Contraseña';
            loginSubtitle.textContent = 'Ingresa tu contraseña actual y la nueva para actualizarla.';
            
            // Hacer requeridos los nuevos campos
            document.getElementById('new-pwd').required = true;
            document.getElementById('new-pwd-confirm').required = true;
      } else {
            // Modo: Iniciar Sesión (Volver a la normalidad)
            labelPwd.textContent = 'Contraseña';
            linkChange.style.display = 'inline';
            
            // Ocultar campos extra con animación
            newFields.style.opacity = '0';
            setTimeout(() => { newFields.style.display = 'none'; }, 300);
            
            // Restaurar textos
            btnSubmit.textContent = 'Iniciar Sesión';
            loginTitle.textContent = 'Bienvenido';
            loginSubtitle.textContent = 'Ingresa tus credenciales para acceder al sistema.';
            
            // Quitar requeridos
            document.getElementById('new-pwd').required = false;
            document.getElementById('new-pwd-confirm').required = false;
            
            // Limpiar valores por seguridad
            document.getElementById('new-pwd').value = '';
            document.getElementById('new-pwd-confirm').value = '';
      }
}

function parsearJsonApi(evt) {
      try {
            // 2. Convertimos la respuesta de texto a un objeto JSON real
            const data = JSON.parse(evt.detail.xhr.response);
            let htmlGenerado = "";

            // 3. Iteramos sobre los resultados y armamos los <li>
            if (data.localidades && data.localidades.length > 0) {
                data.localidades.forEach(loc => {
                    htmlGenerado += `
                        <li class="suggestion-item" onclick="seleccionarCiudad('${loc.nombre}', '${loc.provincia.nombre}')">
                            <span class="city-name">${loc.nombre}</span>
                            <span class="prov-name">${loc.provincia.nombre}</span>
                        </li>
                    `;
                });
            } else {
                // Si la API no encontró la ciudad
                htmlGenerado = `
                    <li class="suggestion-item" style="cursor: default; background: white;">
                        <span class="prov-name">No se encontraron resultados.</span>
                    </li>
                `;
            }

            // 4. ¡La magia! Reemplazamos el JSON crudo por nuestro HTML.
            // Ahora HTMX insertará este HTML directamente en el #suggestionsList
            evt.detail.serverResponse = htmlGenerado;

        } catch (error) {
            console.error("Error al procesar el JSON de Georef:", error);
            evt.detail.serverResponse = ""; // Dejamos vacío en caso de error
        }
}

function seleccionarCiudad(nombre, provincia) {
      searchInput.value = `${nombre}, ${provincia}`;
      document.getElementById('hiddenCity').value = searchInput.value;

      suggestionsList.innerHTML = ''; 
      submitBtn.disabled = false;
}

searchInput.addEventListener('input', () => {
      submitBtn.disabled = true;
});

document.body.addEventListener('htmx:configRequest', function(evt) {
      // Solo modificamos la petición si va dirigida a Georef
      if (evt.detail.path.startsWith("https://apis.datos.gob.ar")) {
            // Eliminamos las cabeceras que causan el bloqueo CORS en APIs públicas
            delete evt.detail.headers['HX-Request'];
            delete evt.detail.headers['HX-Target'];
            delete evt.detail.headers['HX-Current-URL'];
            delete evt.detail.headers['HX-Trigger'];
            delete evt.detail.headers['HX-Trigger-Name'];
      }
})

function selectModal(evt) {
      // Comprobamos que sea el formulario correcto y que el status sea 200 OK
      if (evt.detail.successful) {
            solicitudExitosa(evt);
      } else {
            solicitudErronea();
      }
      // Limpiar el formulario
      evt.detail.elt.reset(); 
}

function solicitudExitosa(evt) {
      // Obtener el email que el usuario ingresó para mostrarlo en el modal
      const email = evt.target.querySelector('input[name="email"]').value;
      document.getElementById('modal-email-target').textContent = email || 'tu correo';

      // Mostrar el modal
      document.getElementById('success-modal').classList.add('active');
}

function solicitudErronea() {
      // Mostrar el modal
      document.getElementById('error-modal').classList.add('active');
}

// Función para cerrar el modal
function closeSuccessModal() {
      document.getElementById('success-modal').classList.remove('active');
}
