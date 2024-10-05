include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/web: run the cmd/web application
.PHONY: run/web
run/web:
	go run ./cmd/web -db-dsn=${DATABASE_URL} -port=4000 -dev

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DATABASE_URL}

## db/migrations/new label=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${label}...'
	migrate create -seq -ext=.sql -dir=./migrations ${label}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DATABASE_URL} up


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...


# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/web: build the cmd/web application
.PHONY: build/web
build/web:
	@echo 'Building cmd/web...'
	go build -ldflags="-s" -o=./bin/web ./cmd/web
