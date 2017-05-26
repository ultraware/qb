all: setup test lint

setup:
	@printf "Getting dependencies ...\n"
	go get -t ./...
	@printf "\n\n"

test:
	@printf "Testing ...\n"
	@go test -cover ./...
	@printf "\n\n"

lint:
	@printf "Running linters ...\n"
	@gometalinter --config .gometalinter.json ./...
	@printf "\n\n"
