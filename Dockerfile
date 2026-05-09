# ETAPA 1: Builder (Imagen pesada con todas las herramientas)
FROM golang:1.24-bookworm AS builder

# Instalar herramientas de generación
RUN go install github.com/a-h/templ/cmd/templ@v0.3.833
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0
RUN go install github.com/Masterminds/squirrel@v1.5.4

WORKDIR /app

# Copiar dependencias de Go para caché
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código y generar/compilar
COPY . .
RUN templ generate && sqlc generate
RUN go build -buildvcs=false -o main .

# ETAPA 2: Final (Imagen liviana para producción)
FROM python:3.13-slim-bookworm

WORKDIR /app

# Instalar solo librerías de ejecución necesarias si tu app de Go llama a scripts de Python
COPY requirements.txt .
RUN pip3 install --no-cache-dir --break-system-packages -r requirements.txt

# COPIAR SOLO EL BINARIO generado en la etapa anterior
COPY --from=builder /app/main .
# Copiar scripts de python o carpetas necesarias para la ejecución (ej: templates)
COPY --from=builder /app/ProcesadoJsons/*.py ./ 
# Si tienes carpetas como 'static' o 'templates', añádelas aquí:
# COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./main"]
