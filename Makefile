download:
	go mod download
	go install github.com/jandelgado/gcov2lcov@latest

test: download
	mkdir -p coverage
	go test -v $(shell go list ./... | grep -v '^github.com/maintc/rustmaps-cli/cmd') -coverprofile=coverage/coverage.out
	gcov2lcov -infile=coverage/coverage.out -outfile=coverage/coverage.lcov
	go tool cover -func=coverage/coverage.out

build: download
	go build -o rustmaps ./