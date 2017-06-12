all:build
	
build:
	gox -osarch="linux/amd64" -output ./bin/market_hqpublish
clean:
	@rm -rf bin
	 
test:
	go test ./go/... -race

out:
	rm hqpublish -rf
	mkdir -p hqpublish/conf.d
	cp -a conf.d/hqpublish_dev.xml hqpublish/conf.d/
	cp -a bin/market_hqpublish hqpublish/
	tar -zcf hqpublish.tar.gz hqpublish/
