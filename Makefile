# Запуск приложения
run:
	go run cmd/auth-medods/main.go

# Миграции
mig_up:
	bash -c 'source .env && migrate -path migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" -verbose up'

mig_down:
	bash -c 'source .env && migrate -path migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" -verbose down'

# Запуск Docker Compose
docker_up:
	docker-compose up -d

docker_down:
	docker-compose down

# Сборка Docker-образа
docker_build:
	docker build -t auth-medods .

# Запуск Docker-контейнера и выполнение тестов
# TODO: fix bug with DNS
# docker_test:
#	docker run --env-file .env.docker -w /usr/local/src auth-medods-app ./scripts/wait-for-it.sh db:5432 --timeout=60 -- go test -v ./...

# Локальные тесты с покрытием
test_coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Локальные тесты с выводом процента покрытия в терминале
test_coverage_terminal:
	go test -cover ./... | grep 'coverage:'

# Локальная версия: Запуск миграций, тестов с покрытием и приложения
local: mig_up test_coverage test_coverage_terminal run

# Docker версия: Сборка образа, запуск контейнера
docker: docker_build docker_up docker_down
