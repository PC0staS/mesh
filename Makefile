.PHONY: build install

build:
	@rm -f build/mesh 2>/dev/null || true
	@go build -o build/mesh .

install: build
	@cp build/mesh /usr/local/bin/