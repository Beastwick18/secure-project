.PHONY: run build test docker docker-run

build:
	go build -o seed-go

run: build
	./seed-go

test:
	go test -v ./...

docker: docker-build docker-run

docker-build:
	docker build -t secure-go .

docker-run:
	docker run -p 8080:8080 secure-go
