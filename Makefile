format:
	@gofmt -s -w .
swagger:
	@swag init -o ./doc/ -g swagger.go -d ./pkg/api/
build:
	@go build
run: build
	@./backupmonitor
win-get:
	@pwsh ./scripts/go-get.ps1
get:
	@./scripts/go-get.sh
