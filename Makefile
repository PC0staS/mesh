.PHONY: build install

build:
	@rm mesh 2>/dev/null || true
	@go build -o mesh .

install: build
	@cp mesh /usr/local/bin/