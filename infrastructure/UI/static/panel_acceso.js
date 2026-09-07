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
                  <input type="text" class="form-control" placeholder="Dr/Dra. Nombre" required>
            </div>
            <div class="vet-input-group">
                  <label class="form-label">Matrícula (MP/MN)</label>
                  <input type="text" class="form-control" placeholder="N° de Matrícula" required>
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
