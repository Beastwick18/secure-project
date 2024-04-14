FROM golang:1.17.2-alpine3.14

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
CMD ["./main"]
