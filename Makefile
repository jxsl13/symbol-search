
.PHONY: build


build:
	@echo "Building..."
	@CGO_ENABLED=0 go build .

install:
	@echo "Installing..."
	@CGO_ENABLED=0 go install .