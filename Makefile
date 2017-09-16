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
	@T=$$(go test -v -short -cover ./...); C=$$?; \
		echo -e "$$T" | grep -vE "^(ok|fail|\?|coverage|PASS)"; echo; \
		echo -e "$$T" | grep -E "^(ok|FAIL)[^$$]"; \
		exit $$C
	@printf "\n\n"

test_all: test_postgres test_mysql test_mssql

test_postgres:
	@TYPE=postgres go test ./tests -v

test_mysql:
	@TYPE=mysql go test ./tests -v

test_mssql:
	@TYPE=mssql go test ./tests -v

lint:
	@printf "Running linters ...\n"
	@gometalinter --config .gometalinter.json ./...
	@printf "\n\n"
