all:build
	
build:
	gox -osarch="linux/amd64" -output ./bin/market_finance
clean:
	@rm -rf bin
	 
test:
	go test ./go/... -race
