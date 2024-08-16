CLI_NAME=celeritas.exe

## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build_cli: builds the command line tool dist directory
dist_cli:
	@go build -o ./dist/${CLI_NAME} ./cmd/cli
