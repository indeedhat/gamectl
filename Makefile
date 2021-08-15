.PHONY: clean
clean:
	rm -rf ./command-center

.PHONY: build
build:
	go build .


.PHONY: run
run:
	go build .
	./command-center
