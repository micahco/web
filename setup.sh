#!/bin/sh

# Install latest migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Execute latest migration
migrate -path=./migrations -database=$DATABASE_URL up
