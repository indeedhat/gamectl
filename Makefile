.PHONY: clean
clean:
	rm -rf ./gamectl

.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./ ./...

.PHONY: run
run: build
	./gamectl
