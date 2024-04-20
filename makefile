.PHONY: run build test docker docker-build docker-run docker-clean
UNUSED_IMAGES := $(shell docker images --filter "dangling=true" -q)

build:
	go build -o seed-go

run: build
	./seed-go ./data/phonebook.db ./data/audit.log

test:
	go test -v ./...

docker: docker-build docker-clean docker-run

docker-build:
	docker build -t secure-go .

docker-run:
	docker run -v ./data:/data -u 1000:1000 -p 8080:8080 secure-go

docker-clean:
	@:$(shell docker rmi $(UNUSED_IMAGES) -f)
