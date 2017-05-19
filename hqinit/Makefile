all:build
	
build:
	gox -osarch="linux/amd64" -output ./bin/market_hqinit
clean:
	@rm -rf bin
	 
test:
	go test ./go/... -race
