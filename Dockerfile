# Builder stage
FROM golang:1.22.3-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . ./

RUN go build -o ./bin/app cmd/auth-medods/main.go

# Runner stage
FROM alpine:latest AS runner

RUN apk add --no-cache ca-certificates postgresql-client

COPY --from=builder /usr/local/src/bin/app /

EXPOSE 8080

CMD ["/app"]
