build:
	@go build -o bin/ecom cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecom

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

# Reload server when file changes are detected.
# go install github.com/cespare/reflex@latest
# See: https://github.com/cespare/reflex
reload:
	@reflex -g '*.go' -s -- sh -c 'make run'