all: setup test lint

setup:
	@printf "Running go generate ...\n"
	@go generate ./...
	@printf "Getting dependencies ...\n"
	@go get -t ./...; true
	@printf "\n\n"

test:
	@printf "Testing ...\n"
	@T=$$(go test -cover ./...); C=$$?; \
		echo -e "$$T" | grep -v "^?"; \
		exit $$C
	@printf "\n\n"

lint:
	@printf "Running linters ...\n"
	@gometalinter --config .gometalinter.json ./...
	@printf "\n\n"

.PHONY:generator
generator:
	@printf "Building qb-generator ...\n"
	@go build -o qb-generator ./generator
