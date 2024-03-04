run-api:
	@echo "Running api..."
	@export $$(cat .env.local) && go run cmd/api/main.go

run-agent:
	@echo "Running agent..."
	@go run cmd/agent/main.go
