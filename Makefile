.PHONY: clean
clean:
	rm -rf ./gamectl

.PHONY: build
build:
	go build -o ./ ./...

.PHONY: run
run: build
	./gamectl
