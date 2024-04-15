.PHONY: run build test docker docker-run

build:
	go build -o seed-go

run: build
	./seed-go ./data/phonebook.db

test:
	go test -v ./...

docker: docker-build docker-run

docker-build:
	docker build -t secure-go .

docker-run:
	docker run -v ./data:/data -u 1000:1000 -p 8080:8080 secure-go
