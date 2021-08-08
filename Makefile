.PHONY: clean
clean:
	rm -rf ./command-center

.PHONY: build
build:
	CGO_ENABLED=0 go build .


.PHONY: run
run:
	CGO_ENABLED=0 go run .
