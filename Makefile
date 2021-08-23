.PHONY: clean
clean:
	rm -rf ./gamectl

.PHONY: build
build:
	go build .

.PHONY: run
run:
	go build .
	./gamectl
