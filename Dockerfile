# Builder stage
FROM golang:1.22.3-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . ./

RUN go build -o ./bin/app cmd/auth-medods/main.go

# Runner stage (с Go для тестов)
FROM golang:1.22.3-alpine AS runner

RUN apk add --no-cache ca-certificates postgresql-client

COPY --from=builder /usr/local/src/bin/app /

# Копируем весь исходный код для запуска тестов
COPY . /usr/local/src/

WORKDIR /usr/local/src/

EXPOSE 8080

CMD ["/app"]
