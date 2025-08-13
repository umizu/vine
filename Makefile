build: 
	@go build -o bin/vine cmd/vine-test/main.go
run: build
	@bin/vine
