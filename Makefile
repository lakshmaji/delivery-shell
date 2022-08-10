build:
	go build main.go

lint:
	golangci-lint run

test:
	go test ./... --coverprofile=coverage.out

coverage:
	go tool cover -html=coverage.out

dev:
	APP_ENVIRONMENT=development go run main.go

start:
	./main


