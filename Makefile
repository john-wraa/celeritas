CLI_NAME=celeritas.exe

## test: runs all tests
test:
	@echo Testing celeritas...
	@go test -v ./...
	@echo Done!

## cover: opens coverage in browser
cover:
	@echo Testing and displaying celeritas cover...
	@go test '-coverprofile=coverage.out' ./...
	@go tool cover '-html=coverage.out'
	@echo Done!

## coverage: displays test coverage
coverage:
	@echo Testing celeritas coverage...
	@go test -cover ./...
	@echo Done!

## build_cli: builds the command line tool dist directory
dist_cli:
	@go build -o ./dist/${CLI_NAME} ./cmd/cli
