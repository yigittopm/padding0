run: build
	@./bin/app

build:
	@go build -o bin/app cmd/padding0/main.go