all:build
	
build:
	gox -osarch="linux/amd64" -output ./bin/hqpost
clean:
	@rm -rf bin
	 
test:
	go test ./go/... -race
