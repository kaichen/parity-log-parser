PROJECTNAME := $(shell basename "$(PWD)")

build:
	@echo "  >  Building binary..."
	@go build -mod=readonly ./...

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo