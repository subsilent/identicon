package web

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "runtime/debug"

    "bitbucket.org/subsilentio/common/middleware"
    "github.com/coocood/freecache"
    "github.com/gin-gonic/gin"
    log "github.com/mgutz/logxi/v1"
    identicon "github.com/subsilent/go-identicon"
)

// Key is used as a seed for then identicon generator
var Key = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

// Options is used to set the server config
type Options struct {
    Logger    log.Logger
    CacheSize int
}

// NewServer creates a new HTTP server from the given config.
func NewServer(opts Options) *gin.Engine {

    // Setup the in-memory cache for icons
    cacheSize := opts.CacheSize * 1024 * 1024
    opts.Logger.Info("Creating cache", "size", cacheSize)
    cache := freecache.NewCache(cacheSize)
    debug.SetGCPercent(10)

    // Create router
    router := gin.New()

    // Setup middleware
    router.Use(gin.Recovery(), middleware.Logger(opts.Logger))

    // Setup icon generation route
    router.GET("/icon/:name", BasicImageRoute(opts.Logger, cache))

    // Setup encoded icon generation route
    router.GET("/encoded/:name", EncodedImageRoute(opts.Logger, cache))

    // Setup random icon generator
    router.GET("/random", RandomImageRoute(opts.Logger))
    return router
}

func cacheHandler(c *gin.Context, l log.Logger, cache *freecache.Cache) (val []byte, err error) {
    name := c.Param("name")
    data := []byte(name)

    // Error check
    if len(data) > 255 {
        return val, fmt.Errorf("Name too long")
    } else if len(data) == 0 {
        return val, fmt.Errorf("Name required")
    }

    // Get the value from the cache
    val, err = cache.Get(data)

    // If the value is not in the cache, generate a new icon and insert it into the cache
    if err != nil || len(val) == 0 {
        icon := identicon.New7x7(Key)

        // Render PNG
        pngdata := icon.Render(data)
        val = pngdata

        // Load value into cache
        if err = cache.Set(data, pngdata, 0); err != nil {
            l.Warn("Failed to load cache")
        }
    }
    return val, nil
}

// BasicImageRoute generates a PNG based on the name provided
func BasicImageRoute(l log.Logger, cache *freecache.Cache) gin.HandlerFunc {
    return func(c *gin.Context) {
        val, err := cacheHandler(c, l, cache)
        if err != nil {
            c.String(400, err.Error())
            return
        }

        // Write response
        c.Writer.Write(val)
        c.Header("Content-Type", "image/png")
    }
}

// EncodedImageRoute encodes the generated image with base64 and returns the string
func EncodedImageRoute(l log.Logger, cache *freecache.Cache) gin.HandlerFunc {
    return func(c *gin.Context) {
        val, err := cacheHandler(c, l, cache)
        if err != nil {
            c.String(400, err.Error())
            return
        }

        encoding := base64.StdEncoding.EncodeToString(val)
        val = []byte("data:image/png;base64," + encoding)

        c.Writer.Write(val)
        c.Header("Content-Type", "text/plain")
    }
}

// RandomImageRoute generates a random image
func RandomImageRoute(l log.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {

        // Create buffer and read random bytes
        data := make([]byte, 16)
        rand.Read(data)

        // Render PNG
        icon := identicon.New7x7(Key)
        pngdata := icon.Render(data)

        // Write response
        c.Writer.Write(pngdata)
        c.Header("Content-Type", "image/png")
    }
}
