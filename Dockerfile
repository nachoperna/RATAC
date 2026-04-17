FROM golang:1.24-bookworm

# Instalación de herramientas (Esto se queda igual)
RUN go install github.com/a-h/templ/cmd/templ@v0.3.833
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0

ENV PATH="/go/bin:${PATH}"
WORKDIR /app

# Solo copiamos los archivos de dependencias para el caché
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

EXPOSE 8080
# El comando ahora lo maneja el docker-compose.yml