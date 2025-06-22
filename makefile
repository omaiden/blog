TEST_FLAGS:=-race -timeout 10s

export TEST_DB_URL=postgres://genie-user:genie-user-pw@localhost:5431/%s?sslmode=disable
export DB_URL=postgres://genie-user:genie-user-pw@localhost:5431/blog?sslmode=disable

run:
	go run cmd/api/main.go
.PHONY: run

test:
	@go test -v ./...
.PHONY: test

test.unit: export SKIP_INTEGRATION=1
test.unit:
	go test ./... ${TEST_FLAGS}
.PHONY: test.unit

vendor:
	go mod tidy
	go mod vendor
.PHONY: vendor
