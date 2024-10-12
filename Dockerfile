# Build the Go app
FROM golang:1.22.0 as builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/padding0/main.go

# Run the Go app
FROM alpine:latest
WORKDIR /
COPY --from=builder . .
CMD ["./main"]