build:
	@go build -o bin/dstore 
run: build
	./bin/dstore

test: 
	go test -v ./...
