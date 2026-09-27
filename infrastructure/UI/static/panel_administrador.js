let currentRowToDelete = null;
let currentRequestCard = null;
let currentRequestAction = null;
let currentCard = null;
offset = 0;
let email = null;

// SVGs para los íconos
const svgQuestion = `<svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>`;

// Función dedicada para Aceptar
function openAcceptModal(btn) {
      currentCard = btn.closest(".request-card");
      email = btn.dataset.email_lab;
      const labName = btn.dataset.nombre_lab;
      document.getElementById('accept-lab-name').textContent = labName;
      document.getElementById('accept-modal').classList.add('active');
}

// Función dedicada para Rechazar
function openRejectModal(btn) {
      currentCard = btn.closest(".request-card");
      email = btn.dataset.email_lab;
      const labName = btn.dataset.nombre_lab;
      document.getElementById('reject-lab-name').textContent = labName;
      document.getElementById('reject-modal').classList.add('active');
}

function processAction() {
      setTimeout(() => { }, 400)
      // Animación fluida para "hacer desaparecer" la tarjeta
      currentCard.style.transition = 'all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275)';
      currentCard.style.opacity = '0';
      currentCard.style.transform = 'scale(0.9)';
      
      // Si quieres que el espacio colapse
      setTimeout(() => {
            currentCard.style.display = 'none';
      }, 400);
}

function closeModals() {
      document.querySelector('.modal-overlay.active').classList.remove('active');
}

function habilitarAnterior() {
      document.getElementById('btn-anterior').disabled = false;
}
