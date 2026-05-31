# Variable para acortar los comandos de ejecución en Docker
DB_URL = postgres://admin:password@db:5432/RATAC_DB?sslmode=disable

MIGRATE_DOCKER = docker run --rm -v "$(PWD)/$(MIGRATIONS_DIR):/migrations:z" --network ratac_default migrate/migrate

DOCKER_EXEC = docker compose exec app
MIGRATIONS_DIR=./DB/migrations
# Variables para herramientas de generación efímeras (con permisos corregidos)
SQLC_DOCKER = docker run --rm -u $(shell id -u):$(shell id -g) -v "$(PWD):/src:z" -w /src sqlc/sqlc:1.28.0
TEMPL_DOCKER = docker run --rm -u $(shell id -u):$(shell id -g) -v "$(PWD):/app:z" -w /app ghcr.io/a-h/templ:v0.3.833

# Fuerza el build para aplicar cambios en el Dockerfile/requirements
bdocker:
	docker compose up -d --build

# Levanta los contenedores 
udocker:
	docker compose up -d db

up-appdocker:
	docker compose up -d app

down-appdocker:
	docker compose down app

# Baja los contenedores y borra los volúmenes (limpieza total de la DB)
dvdocker:
	docker compose down -v 

# Baja los contenedores normalmente
ddocker:
	docker compose down 

# Crear una nueva migración
# Uso: make migrate-create name=nombre_migracion
migrate-up:
	$(MIGRATE_DOCKER) -path=/migrations -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE_DOCKER) -path=/migrations -database "$(DB_URL)" down

migrate-create:
	$(MIGRATE_DOCKER) create -ext sql -dir /migrations -seq $(name)

# Revertir todas las migraciones
migrate-reset:
	$(MIGRATE_DOCKER) -path=/migrations -database "$(DB_URL)" down -all

# Borra copia de imagenes antiguas
clean-images:
	docker system prune -f

# PROCESAMIENTO: Ahora corre dentro del contenedor usando las dependencias de Python instaladas
procesarjsons: udocker wait up-appdocker wait
	$(DOCKER_EXEC) ./ProcesadoJsons/eliminar_duplicados.sh
	$(DOCKER_EXEC) python3 ./ProcesadoJsons/diag_to_json.py
	$(DOCKER_EXEC) python3 ./ProcesadoJsons/PDF_to_json.py  
	$(DOCKER_EXEC) go run ./ProcesadoJsons/json_to_bd.go 
	$(MAKE) down-appdocker

# Ejecución local (por si querés probar algo fuera de Docker, requiere dependencias locales)
run:
	go run ./main.go 

wait:
	@sleep 2

dependencias:
	$(SQLC_DOCKER) generate
	$(TEMPL_DOCKER) generate

# Acceso directo a la terminal de la base de datos PostgreSQL
sql-directo:
	docker exec -it RATAC_db psql -U admin -d RATAC_DB 

# Útil para ver qué está pasando dentro del contenedor en tiempo real
logs:
	docker compose logs -f app

# Alias para levantar todo el entorno
server: udocker wait dependencias run

emptyJSONS:
	rm -rf JSONS/

.PHONY: run udocker dvdocker ddocker procesarjsons sql-directo logs server bdocker wait clean-images emptyJSONS dependencias
