build:
	go build -o=./bin/web ./cmd/web

run:
	go run ./cmd/web -dev -port=4000
