
COVERAGE_THRESHOLD=80
.PHONY: run swag test swagrun

run:
	go run ./command/api/main.go

swag:
	swag init -g ./command/api/main.go -o ./docs

test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	python python_scripts/check_coverage.py --threshold $(COVERAGE_THRESHOLD)
	go tool cover -html=coverage.out -o coverage.html

swag-run: swag run
