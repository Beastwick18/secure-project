FROM golang:1.19.1-alpine3.16

RUN apk add --no-cache libc-dev
RUN apk add --no-cache gcc
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go test -v ./...
RUN go build -o main .

EXPOSE 8080

# Run the executable
CMD ["./main", "/users.db"]
