
.PHONY: run swag test swagrun

run:
	go run ./command/api/main.go

swag:
	swag init -g ./command/api/main.go -o ./docs

test:
	go test ./... -coverpkg=./... -covermode=atomic -p 1 -coverprofile=coverage.tmp
	python filter_coverage.py
	go tool cover -html=coverage.out -o coverage.html

swagrun: swag run
