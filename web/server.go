package web

import (
    "encoding/base64"
    "runtime/debug"
    "time"

    "github.com/coocood/freecache"
    identicon "github.com/dgryski/go-identicon"
    "github.com/gin-gonic/gin"
    log "github.com/mgutz/logxi/v1"
)

var Key = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

type Options struct {
    Logger    log.Logger
    CacheSize int
}

func NewServer(opts Options) *gin.Engine {

    // options
    logger := opts.Logger

    cacheSize := opts.CacheSize * 1024 * 1024
    logger.Info("Creating cache", "size", cacheSize)
    cache := freecache.NewCache(cacheSize)
    debug.SetGCPercent(10)

    // create router
    router := gin.New()
    router.Use(gin.Recovery(), Logger(logger))
    router.GET("/icon/:name", BasicImageRoute(logger))
    router.GET("/encoded/:name", EncodedImageRoute(logger, cache))
    return router
}

func Logger(l log.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        start := time.Now()
        path := c.Request.URL.Path

        // Process request
        c.Next()

        // Stop timer
        end := time.Now()
        latency := end.Sub(start)

        // clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()
        comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

        l.Info("Request", "method", method, "status", statusCode, "latency", latency, "path", path, "comment", comment)
    }
}

func BasicImageRoute(l log.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        icon := identicon.New7x7(Key)
        name := c.Param("name")

        l.Info("Creating identicon", "name", name)

        data := []byte(name)
        pngdata := icon.Render(data)

        c.Writer.Write(pngdata)
        c.Header("Content-Type", "image/png")
    }
}

func EncodedImageRoute(l log.Logger, cache *freecache.Cache) gin.HandlerFunc {
    return func(c *gin.Context) {

        icon := identicon.New7x7(Key)
        name := c.Param("name")
        data := []byte(name)
        if len(data) > 255 {
            c.String(400, "Name too long")
            return
        }

        val, err := cache.Get(data)
        if err != nil || len(val) == 0 {
            l.Info("Encoding identicon", "name", name)
            pngdata := icon.Render(data)
            encoding := base64.StdEncoding.EncodeToString(pngdata)
            val = []byte("data:image/png;base64," + encoding)

            l.Info("Loading cache")
            if err = cache.Set(data, val, 0); err != nil {
                l.Info("Failed to load cache")
            }
        }

        c.Writer.Write(val)
        c.Header("Content-Type", "text/plain")
    }
}
