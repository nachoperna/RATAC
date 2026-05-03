# Variable para acortar los comandos de ejecución en Docker
DOCKER_EXEC = docker compose exec app

# Levanta los contenedores y fuerza el build para aplicar cambios en el Dockerfile/requirements
udocker:
	docker compose up -d --build

# Baja los contenedores y borra los volúmenes (limpieza total de la DB)
dvdocker:
	docker compose down -v 

# Baja los contenedores normalmente
ddocker:
	docker compose down 

# PROCESAMIENTO: Ahora corre dentro del contenedor usando las dependencias de Python instaladas
procesarjsons:
	$(DOCKER_EXEC) python3 ./ProcesadoJsons/diag_to_json.py
	$(DOCKER_EXEC) python3 ./ProcesadoJsons/PDF_to_json.py  
	$(DOCKER_EXEC) go run ./ProcesadoJsons/json_to_bd.go 

# Ejecución local (por si querés probar algo fuera de Docker, requiere dependencias locales)
run:
	go run ./main.go 

# Acceso directo a la terminal de la base de datos PostgreSQL
sql-directo:
	docker exec -it RATAC_db psql -U admin -d RATAC_DB 

# Útil para ver qué está pasando dentro del contenedor en tiempo real
logs:
	docker compose logs -f app

# Alias para levantar todo el entorno
server: udocker

emptyJSONS:
	cd JSONS && sudo rm -rf * && cd ..

.PHONY: run udocker dvdocker ddocker procesarjsons sql-directo logs server