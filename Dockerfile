# Builder stage
FROM golang:1.22.3-alpine AS builder

WORKDIR /usr/local/src

# Копируем go.mod и go.sum для управления зависимостями
COPY ["go.mod", "go.sum", "./"]

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . ./

# Сборка Go-приложения
RUN go build -o ./bin/app cmd/auth-medods/main.go

# Runner stage (с Go для тестов)
FROM golang:1.22.3-alpine AS runner

# Установка необходимых зависимостей
RUN apk add --no-cache ca-certificates postgresql-client bash

# Копируем скомпилированное приложение из builder stage
COPY --from=builder /usr/local/src/bin/app /

# Копируем весь исходный код для запуска тестов и работы приложения
COPY . /usr/local/src/

# Копируем скрипт wait-for-it.sh
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /usr/local/bin/wait-for-it.sh

# Делаем скрипт исполняемым
RUN chmod +x /usr/local/bin/wait-for-it.sh

WORKDIR /usr/local/src/

# Открываем порт для приложения
EXPOSE 8080

# Стартовое командное приложение
CMD ["/app"]
