
COVERAGE_THRESHOLD=80

COVERAGE_EXCLUDE = main\.go|api\.go|redis_client_mock\.go|test|docs|_func_demo|config|internal/service/mocks/|internal/service/services_mocks/|pkg/logger/|pkg/utils_log/

.PHONY: run swag test swag-run

run:
	go run ./command/api/main.go

swag:
	swag init -g ./command/api/main.go -o ./docs

test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -Ev '$(COVERAGE_EXCLUDE)' coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep 'total:' | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$total%"; \
	if [ "$$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l)" -eq 1 ]; then \
		echo "ERROR: Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "SUCCESS: Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

swag-run: swag run
