package main

import (
	"flag"
	"net/http"
	"runtime"

	"os"

	"github.com/facebookgo/grace/gracehttp"
	log "github.com/mgutz/logxi/v1"
	"github.com/subsilent/identicon/web"
)

var (
	address = flag.String("addr", ":3000", "Address to bind to")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	// Create logger
	writer := log.NewConcurrentWriter(os.Stdout)
	logger := log.NewLogger(writer, "identicon")

	// Start web server
	router := web.NewServer(logger)

	gracehttp.Serve(
		&http.Server{Addr: *address, Handler: router},
	)
}
