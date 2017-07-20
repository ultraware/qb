all: setup test lint

setup:
	@go install ./qb-generator
	@printf "Running go generate ...\n"
	@go generate ./...
	@printf "Getting dependencies ...\n"
	@go get -t ./...; true
	@printf "\n\n"

test:
	@printf "Testing ...\n"
	@T=$$(go test -short -cover ./... $$VERBOSE); C=$$?; \
		echo -e "$$T" | grep -v "^?" | grep -v " 0.0%"; \
		exit $$C
	@printf "\n\n"

test_all: test_postgres test_mysql test_mssql

test_postgres:
	@TYPE=postgres go test ./tests $$VERBOSE

test_mysql:
	@TYPE=mysql go test ./tests $$VERBOSE

test_mssql:
	@TYPE=mssql go test ./tests $$VERBOSE

lint:
	@printf "Running linters ...\n"
	@gometalinter --config .gometalinter.json ./...
	@printf "\n\n"
