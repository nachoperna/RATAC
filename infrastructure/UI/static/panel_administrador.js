let currentRowToDelete = null;
let currentRequestCard = null;
let currentRequestAction = null;
offset = 0;
// ==========================================
// Lógica de Eliminación de Filas
// ==========================================
function openDeleteModal(btn) {
      currentRowToDelete = btn.closest('tr');
      const protocolo = currentRowToDelete.querySelector('.cell-protocol').innerHTML;
      const textElement = document.getElementById('delete-modal-text');
      textElement.innerHTML = `¿Está seguro que desea eliminar el diagnóstico con protocolo <strong>${protocolo}</strong> de la base de datos? Esta acción no se puede deshacer.`;

      document.getElementById('delete-modal').classList.add('active');
}

function confirmDelete(btn) {
      currentRowToDelete = btn.closest('tr');
      // Pequeña animación de desvanecimiento
      currentRowToDelete.style.transition = 'opacity 0.3s ease';
      currentRowToDelete.style.opacity = '0';
      setTimeout(() => {
            currentRowToDelete.remove();
            currentRowToDelete = null;
            closeModals();
      }, 300);
}

// ==========================================
// Lógica de Solicitudes (Aceptar/Rechazar)
// ==========================================
function openRequestModal(btn, action, labName) {
      currentRequestCard = btn.closest('article');
      currentRequestAction = action;

      const titleElement = document.getElementById('request-modal-title');
      const confirmBtn = document.getElementById('btn-confirm-request');

      if(action === 'aceptar') {
            titleElement.textContent = `Aceptar: ${labName}`;
            confirmBtn.className = 'btn btn-primary';
            confirmBtn.textContent = 'Aceptar y Habilitar';
      } else {
            titleElement.textContent = `Rechazar: ${labName}`;
            confirmBtn.className = 'btn btn-danger';
            confirmBtn.textContent = 'Rechazar Solicitud';
      }

      document.querySelector('.modal-input').value = ''; // limpiar textarea
      document.getElementById('request-modal').classList.add('active');
}

function submitRequest() {
      if(currentRequestCard) {
            // Animación para eliminar la tarjeta aceptada/rechazada
            currentRequestCard.style.transition = 'all 0.3s ease';
            currentRequestCard.style.opacity = '0';
            currentRequestCard.style.transform = 'scale(0.95)';
            setTimeout(() => {
                  currentRequestCard.remove();
                  currentRequestCard = null;
                  currentRequestAction = null;
                  closeModals();
            }, 300);
      }
}

// ==========================================
// Utilidades Generales
// ==========================================
function closeModals(event) {
      document.querySelectorAll('.modal-overlay').forEach(modal => {
            modal.classList.remove('active');
      });
      currentRowToDelete = null;
      currentRequestCard = null;
      currentRequestAction = null;
}

// Cerrar al hacer click afuera de la caja del modal
document.querySelectorAll('.modal-overlay').forEach(overlay => {
      overlay.addEventListener('click', (e) => {
            if(e.target === overlay) {
                  closeModals(e);
            }
      });
});

function habilitarAnterior() {
      document.getElementById('btn-anterior').disabled = false;
}
