all:build
	
build:
	gox -osarch="linux/amd64" -output ./bin/market_hqpost
clean:
	@rm -rf bin
	 
test:
	go test ./go/... -race
