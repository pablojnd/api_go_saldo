# Etapa de compilación
FROM golang:1.20 as build

WORKDIR /app

# Copiar los archivos de dependencias primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o /api_server .

# Etapa final con una imagen mínima
FROM alpine:latest

# Añadir certificados SSL y zona horaria
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copiar el binario compilado
COPY --from=build /api_server /api_server

# Copiar el archivo .env (opcional)
# COPY --from=build /app/.env /.env

# Exponer el puerto configurado en la aplicación
EXPOSE 8080

# Ejecutar la aplicación
CMD ["/api_server"]
