proj = identicon
binary = $(proj)

build: **/*.go
	@go build .

test: **/*.go
	@go test

deps:
	@go get .

clean:
	@rm $(binary)

run: export LOGXI=*=WRN
run: export LOGXI_COLORS= *=black,key=black+h,message=blue,TRC,DBG,WRN=red+h,INF=green,ERR=red+h,maxcol=1000
run: export LOGXI_FORMAT= happy,t=2006-01-02 15:04:05.000000
run: build
	./$(binary)

.PHONY: build deps clean run
