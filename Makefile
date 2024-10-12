run: build
	@./bin/app

build:
	@go build -o bin/padding0 ./main.go