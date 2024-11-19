build:
	@go build -o bin/scratch cmd/main.go

test:
	@go test ./...

run: build
	@./bin/scratch

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down
