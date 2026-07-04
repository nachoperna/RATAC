# ETAPA 1: Builder (Imagen pesada con todas las herramientas)
FROM golang:1.24-bookworm AS builder

# Instalar herramientas de generación
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1001
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0

# Instalar solo librerías de ejecución necesarias si tu app de Go llama a scripts de Python
# 1. Instalar pip3 a nivel de sistema operativo
RUN apt-get update && apt-get install -y python3-pip
COPY requirements.txt .
RUN pip3 install --no-cache-dir --break-system-packages -r requirements.txt

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

# COPIAR SOLO EL BINARIO generado en la etapa anterior
COPY --from=builder /app/main .
# Copiar scripts de python o carpetas necesarias para la ejecución (ej: templates)
COPY --from=builder /app/ProcesadoJsons ./ProcesadoJsons
COPY --from=builder /app/infrastructure/UI/static ./static

EXPOSE 8080

CMD ["./main"]
