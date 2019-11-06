.PHONY: all
all: setup test lint

.PHONY: setup
setup:
	@go install ./qb-generator
	@printf "Running go generate ...\n"
	@go generate ./...
	@printf "Getting dependencies ...\n"
	@go get -t ./...
	@printf "\n\n"

.PHONY: test
test:
	@printf "Testing ...\n"
	@go test -short -cover -coverprofile unit.out ./...
	@printf "\n\n"

.PHONY: test_all
test_all: test_postgres test_mysql test_mssql

.PHONY: test_postgres
test_postgres:
	@TYPE=postgres go test -cover -coverprofile postgres.out -coverpkg ./... ./internal/tests -v

.PHONY: test_mysql
test_mysql:
	@TYPE=mysql go test -cover -coverprofile mysql.out -coverpkg ./... ./internal/tests -v

.PHONY: test_mssql
test_mssql:
	@TYPE=mssql go test -cover -coverprofile mssql.out -coverpkg ./... ./internal/tests -v

.PHONY: lint
lint:
	@printf "Running linters ...\n"
	@golangci-lint run
	@printf "\n\n"
