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

run: export LOGXI=*=INF
run: export LOGXI_COLORS= *=black,key=black+h,message=blue,TRC,DBG,WRN=red+h,INF=green,ERR=red+h,maxcol=1000
run: export LOGXI_FORMAT= happy,t=2006-01-02 15:04:05.000000
run: build
	./$(binary)

docker: export GOOS=linux
docker: export CGO_ENABLED=0
docker: export GOARCH=amd64
docker: 
	@godep go build -a -installsuffix cgo -o main .
	@docker build -t identicon -f Dockerfile.scratch .

.PHONY: build deps clean run
