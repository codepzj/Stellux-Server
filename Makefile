.PHONY: wire fmt

wire:
	wire ./...

fmt:
	go fmt ./...
