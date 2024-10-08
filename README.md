# Test Task BackDev

Тестовое задание на позицию Junior Backend Developer

## Используемые технологии

- Go
- JWT
- PostgreSQL
- Docker

## Задание

Написать часть сервиса аутентификации.

### Маршруты

1. **Маршрут для получения пары Access, Refresh токенов**
   - Выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса.

2. **Маршрут для выполнения Refresh операции**
   - Выполняет Refresh операцию на пару Access, Refresh токенов.

## Требования

- **Access токен**: тип JWT, алгоритм SHA512, хранить в базе строго запрещено.
- **Refresh токен**: тип произвольный, формат передачи base64, хранится в базе исключительно в виде bcrypt хеша, должен быть защищен от изменения на стороне клиента и попыток повторного использования.
- Access и Refresh токены обоюдно связаны. Refresh операцию для Access токена можно выполнить только тем Refresh токеном, который был выдан вместе с ним.
- Payload токенов должен содержать сведения об IP адресе клиента, которому он был выдан. В случае, если IP адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные).

## Результат

Результат выполнения задания нужно предоставить в виде исходного кода на Github. Будет плюсом, если получится использовать Docker и покрыть код тестами.

## Переменные окружения

Приложение использует различные переменные окружения для настройки базы данных и других параметров. В зависимости от среды, вы можете использовать разные значения для локальной разработки и для Docker-контейнеров.

### Локальная среда (local)

Для локальной разработки используйте следующие переменные окружения:

```plaintext
# Переменные для подключения к базе данных (локальная среда)
DB_USER=postgres
DB_PASSWORD=mypass
DB_HOST=localhost
DB_PORT=5432
DB_NAME=postgres
DB_CONN_STR=postgres://postgres:mypass@localhost:5432/postgres?sslmode=disable

# Secret key и другие настройки
OWN_KEY=volchok96
TOKEN_TTL=30m
```

### Среда Docker

Для запуска в контейнерах Docker используются следующие переменные окружения:

```plaintext
# Переменные для подключения к базе данных (Docker)
DB_USER=postgres
DB_PASSWORD=mypass
DB_HOST=db
DB_PORT=5432
DB_NAME=postgres
DB_CONN_STR=postgres://postgres:mypass@db:5432/postgres?sslmode=disable

# Secret key и другие настройки
OWN_KEY=volchok96
TOKEN_TTL=30m
```

### Настройка переменных окружения

1. **Локальная среда:** Для запуска приложения в локальной среде создайте файл `.env` в корне проекта и добавьте в него переменные окружения, указанные выше для локальной среды.

2. **Среда Docker:** Для Docker создайте файл `.env.docker` и добавьте в него переменные окружения, указанные выше для Docker. Убедитесь, что ваш `docker-compose.yml` файл загружает этот файл с помощью параметра `--env-file`.

```yaml
version: '3'
services:
  app:
    env_file:
      - .env.docker
    ...
```

### Переменные

- **`DB_USER`**: Имя пользователя для базы данных.
- **`DB_PASSWORD`**: Пароль пользователя базы данных.
- **`DB_HOST`**: Хост, на котором находится база данных (локально — `localhost`, в Docker — `db`).
- **`DB_PORT`**: Порт для подключения к базе данных (обычно 5432).
- **`DB_NAME`**: Имя базы данных.
- **`DB_CONN_STR`**: Полная строка подключения к базе данных.
- **`OWN_KEY`**: Секретный ключ для подписи JWT токенов.
- **`TOKEN_TTL`**: Время жизни токенов (например, `30m` для 30 минут).

Эти переменные можно изменить в зависимости от требований вашей среды.

## Установка

1. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/volchok96/auth-medods.git
   cd auth-medods
   ```

2. Установите зависимости:
   ```sh
   go mod tidy
   ```

3. Запустите приложение:
   ```sh
   go run main.go
   ```

## Использование

### Получение токенов

```sh
GET /access?guid=GUID
```

### Refresh токенов

```sh
POST /refresh
Body:
{
  "GUID": "GUID",
  "refresh_token": "refresh_token"
}
```

## Логика email уведомлений при смене IP пользователя

При выполнении операции Refresh токенов, приложение проверяет, изменился ли IP адрес пользователя. Если IP адрес изменился, приложение отправляет email уведомление пользователю. Это делается для обеспечения безопасности и предотвращения несанкционированного доступа.

## Тестирование

В приложении подключено:
- юнит-тестирование (****покрытие: 50.3%****)
- интеграционное тестирование (****покрытие: 22.8%****).

Для запуска тестов используйте команду:
```sh
go test ./...
```

## Docker

Для запуска приложения в Docker используйте команду:
```sh
docker-compose up --build
```

## Makefile

### Запуск приложения

```sh
make run
```

### Запуск Docker Compose

```sh
make docker_up
make docker_down
```

### Сборка Docker-образа

```sh
make docker_build
```

### Локальные тесты с покрытием

```sh
make test_coverage
make test_coverage_terminal
```

### Docker тесты 

```sh
make docker_test
```

### Локальная версия: Запуск тестов с покрытием и приложения

```sh
make local
```

### Docker версия: Сборка образа, запуск контейнера, запуск тестов

```sh
make docker
```

## Используемые библиотеки

- **github.com/dgrijalva/jwt-go**: Библиотека для работы с JWT токенами.
- **github.com/go-chi/chi/v5**: HTTP-router для Go, используемый для создания маршрутов.
- **github.com/google/uuid**: Библиотека для генерации UUID.
- **github.com/lib/pq**: Драйвер для работы с PostgreSQL.
- **github.com/rs/zerolog**: Библиотека для логирования.
- **github.com/stretchr/testify**: Библиотека для тестирования.
- **golang.org/x/crypto**: Библиотека для криптографических операций.
- **gopkg.in/mail.v2**: Библиотека для отправки электронной почты.

## Контакты

Если у вас есть вопросы или предложения, пожалуйста, свяжитесь со мной:
- Email: kzakharova96@yandex.com
- TG: https://t.me/volchok_96