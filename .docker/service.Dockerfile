# Stage 1: Build
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates bash

WORKDIR /app

# Copia solo go.mod, senza go.sum
COPY go.mod ./

# Risolve le dipendenze direttamente durante il build
RUN go mod tidy
RUN go mod download

# Copia tutto il codice
COPY . .

# Build binario
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./src/cmd/app

# Stage 2: Runtime minimal
FROM scratch

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 3000
CMD ["./app"]