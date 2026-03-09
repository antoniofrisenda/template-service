FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./src/cmd/app

FROM scratch

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 3000

CMD ["./app"]