# Stage 1: Build
FROM golang:1.26-alpine AS builder

# Installa strumenti necessari
RUN apk add --no-cache git ca-certificates bash

WORKDIR /app

# Copia solo i file dei moduli per sfruttare la cache Docker
COPY go.mod go.sum ./

# Scarica tutte le dipendenze, cache-friendly
RUN go mod download
RUN go mod tidy

# Copia tutto il codice sorgente
COPY . .

# Build binario Go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./src/cmd/app

# Stage 2: Minimal image
FROM alpine:3.18 AS runtime

WORKDIR /app

# Copia il binario compilato dal builder
COPY --from=builder /app/app .

# Esponi porta
EXPOSE 3000

# Comando di default
CMD ["./app"]