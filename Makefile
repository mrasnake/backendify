build: ## Build server executable.
	go build -o .

run: build ## Build and run server executable
	./backendify

test: build ## Run CLI help flag
	go test