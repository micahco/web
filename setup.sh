#!/bin/sh

# Install go CLI tools
#   NOTE:
# It's important to install these in a Dev Container "postCreateCommand" script,
# instead of inside the Dockerfile. If these were executed within the Dockerfile,
# then they would be executed as root and mess up file permissions.
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# Execute latest migration
migrate -path=./migrations -database=$DATABASE_URL up
