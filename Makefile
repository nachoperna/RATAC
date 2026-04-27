udocker:
	docker compose up -d

dvdocker:
	docker compose down -v

ddocker:
	docker compose down

procesarjsons:
	python3 ./ProcesadoJsons/diag_to_json.py
	go run ./ProcesadoJsons/json_to_bd.go

run:
	go run ./main.go

sql-directo:
	docker exec -it RATAC_db psql -U admin -d RATAC_DB

server: udocker run

.PHONY: run udocker dvdocker ddocker
