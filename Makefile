GO		:= $(shell which go)

.PHONY: test
test:
	$(GO) test -v ./...

.PHONY: run
run:
	docker compose up