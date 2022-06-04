.PHONY: clean
clean:
	rm -rf ./gamectl

.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./ ./...

.PHONY: run
run: build
	./gamectl

.PONY: deploy
deploy:
	ssh root@mc.phpmatt.com systemctl stop gamectl
	rsync --update -aRz --progress gamectl web/ root@mc.phpmatt.com:/opt/gamectl/
	ssh root@mc.phpmatt.com systemctl start gamectl
