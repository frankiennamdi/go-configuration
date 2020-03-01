.PHONY: test
.DEFAULT_GOAL := test

dependencies:
	go mod vendor; go mod tidy

test: dependencies
	go test ./... -v
