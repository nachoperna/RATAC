FROM golang:1.24-bookworm

# 1. Instalamos Python y pip en el contenedor
RUN apt-get update && apt-get install -y python3 python3-pip

# Instalación de herramientas de Go 
RUN go install github.com/a-h/templ/cmd/templ@v0.3.833
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0

ENV PATH="/go/bin:${PATH}"
WORKDIR /app

# 2. Copiamos el requirements.txt ANTES del código para aprovechar el caché
COPY requirements.txt ./
# Instalamos las dependencias de Python
RUN pip3 install --break-system-packages -r requirements.txt

# Solo copiamos los archivos de dependencias de Go para el caché
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código (incluyendo tus .py)
COPY . .

EXPOSE 8080
# El comando ahora lo maneja el docker-compose.yml