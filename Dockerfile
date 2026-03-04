FROM golang:1.25-alpine AS Go

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./src/cmd/app

FROM scratch

WORKDIR /app

COPY --from=Go /app/app .

EXPOSE 3000

CMD ["./app"]