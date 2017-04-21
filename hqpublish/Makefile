all:build
	
build:
	gox -osarch="linux/amd64" -output ./bin/market_hqpublish
clean:
	@rm -rf bin
	 
test:
	go test ./go/... -race
