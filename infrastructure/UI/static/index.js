let filterCount = 0;
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
      nuevoFiltro.insertBefore(operador_logico, nuevoFiltro.firstChild);
      const filtros = document.querySelector('.filtros');
      filtros.insertBefore(nuevoFiltro, filtros.querySelector('button'));
}

function addInputOr(elemento) {
  const filtro = elemento.closest('.filtro');
  const inputs = filtro.querySelectorAll('input.filtro-valor');
  const lastInput = inputs[inputs.length - 1];
  
  const nuevoInput = document.createElement('input');
  nuevoInput.type = 'text';
  nuevoInput.className = 'filtro-valor';
  
  const btnAgregar = document.createElement('button');
  btnAgregar.type = 'button';
  btnAgregar.textContent = '+';
  btnAgregar.onclick = function() { addInputOr(this); };
  btnAgregar.style.marginLeft = '5px';
  
  lastInput.after(nuevoInput, btnAgregar);
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

function removeFilter(btn) {
      btn.parentElement.remove();
}
