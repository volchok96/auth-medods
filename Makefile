run:
	go run cmd/auth-medods/main.go

mig_up:
	migrate -path migrations -database "postgres://postgres:mypass@localhost:5432/postgres?sslmode=disable" -verbose up

mig_down:
	migrate -path migrations -database "postgres://postgres:mypass@localhost:5432/postgres?sslmode=disable" -verbose down
