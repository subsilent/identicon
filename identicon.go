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
    address   = flag.String("addr", ":3000", "Address to bind to")
    cachesize = flag.Int("cachesize", 128, "Size of image cache in MB")
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    flag.Parse()

    // Create logger
    writer := log.NewConcurrentWriter(os.Stdout)
    logger := log.NewLogger(writer, "identicon")
    // logger.SetLevel(log.LevelCritical)

    // Start web server
    router := web.NewServer(web.Options{
        Logger:    logger,
        CacheSize: *cachesize,
    })

    gracehttp.Serve(
        &http.Server{Addr: *address, Handler: router},
    )
}
