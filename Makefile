GOBIN = build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run
DEV ?= false  # Default to false if not set in the environment



build:
	@echo "Building binder."
	go build -o $(GOBIN)/abibinder -buildvcs=false ./cmd/abibinder
	@echo "Done building binder."
	@echo "Run \"$(GOBIN)/abibinder\" to launch abi binder."


.PHONY: build