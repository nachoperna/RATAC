function addFiltro() {
      const filtroPrimero = document.querySelector('.filtro');
      
      const nuevoFiltro = filtroPrimero.cloneNode(true);
      nuevoFiltro.querySelectorAll('input').forEach(input => {
            if (input.type === 'checkbox') input.checked = false;
            else input.value = '';
      });
      const operador_logico = document.createElement('select');
      operador_logico.className = 'operador-logico';
      operador_logico.innerHTML = `
            <option value="AND">AND</option>
            <option value="OR">OR</option>
      `
      const nuevo_header = nuevoFiltro.querySelector('.filtro-header');
      nuevo_header.insertBefore(operador_logico, nuevo_header.firstChild);
      const filtros = document.querySelector('.filtros');
      filtros.insertBefore(nuevoFiltro, filtros.querySelector('.filtros-footer'));
}

function addInputOr(elemento) {
      const esCheckbox = elemento.type === 'checkbox';

      const filtro = elemento.closest('.filtro');
      const inputs = filtro.querySelectorAll('input.filtro-valor');
      if (esCheckbox){ // significa que el usuario quito la seleccion multiple y quitamos todos los inputs y botones extras
            filtro.querySelector('.filtro-add-valor-container').classList.toggle('hidden')
            if (!elemento.checked){
                  const botones = filtro.querySelector('.filtro-add-valor');
                  for (let i = inputs.length - 1; i > 0; i--) {
                        inputs[i].remove();
                  }
                  botones.forEach(btn => btn.remove());
                  return
            }
      }
      const lastInput = inputs[inputs.length - 1];
      
      const nuevoInput = document.createElement('input');
      nuevoInput.type = 'text';
      nuevoInput.className = 'filtro-valor';
      nuevoInput.placeholder = 'Otro valor...';
      
      lastInput.after(nuevoInput);
}

function obtenerFiltros() {
      const filas = document.querySelectorAll('.filtro');
      const arrayFiltros = [];
      const arrayValores = [];

      filas.forEach((fila, index) => {
            // El primer elemento no tiene lógica (AND/OR), los demás sí
            const logicaSelect = fila.querySelector('.operador-logico');
            const logica = (index > 0 && logicaSelect) ? logicaSelect.value : "";
            const campo = fila.querySelector('.filtro-selector').value;
            const operador = fila.querySelector('.operador').value;
            const valoresNodes = fila.querySelectorAll('.filtro-valor');
            const valores = Array.from(valoresNodes).map(input => input.value);
            const isNot = fila.querySelector('.filtro-not').checked;
            const isMultiple = fila.querySelector('.filtro-or').checked;

            valores.forEach(valor => {arrayValores.push(valor)})

            // Añadimos el objeto al array. Esto garantiza que se mantenga el orden estricto de arriba hacia abajo.
            arrayFiltros.push({
                  logica: logica,
                  campo: campo,
                  operador: operador,
                  valores: valores,
                  not: isNot,
                  multiple: isMultiple
            });
      });
      arrayFiltros.forEach(filtro => {
      console.log("Lógica:", filtro.logica);
      console.log("Campo:", filtro.campo);
      console.log("Operador:", filtro.operador);
      console.log("Valores:", filtro.valores);
      console.log("Not:", filtro.not);
      console.log("Multiple:", filtro.multiple);
      console.log("---");
      });
      console.log("Enviando payload:", arrayFiltros);
      return arrayFiltros;
}

function removeFiltro(btn) {
      if (document.querySelectorAll('.filtro').length == 1) {
          toggleFiltros();  
      }else{
            btn.closest('.filtro').remove();
      }
}

function toggleFiltros(){
      document.querySelector('.filtros').classList.toggle('hidden');
}

document.addEventListener('DOMContentLoaded', function() {
      document.querySelector('.filtros').addEventListener('change', function(e) {
      if (e.target.classList.contains('filtro-selector')) {
            const es_edad = e.target.value === 'Edad';
            const operador_select = e.target.closest('.filtro-header').querySelector('.operador');
            Array.from(operador_select.options).forEach(option => {
                  if (option.value !== 'igual') {
                        option.classList.toggle('hidden', !es_edad);
                  }
            });
      }
      });
});
