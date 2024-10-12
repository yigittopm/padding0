run: build
	@./bin/padding0

build:
	@go build -o bin/padding0 ./main.go