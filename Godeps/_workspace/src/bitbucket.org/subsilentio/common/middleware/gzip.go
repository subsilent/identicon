package middleware

import (
    "compress/gzip"

    "github.com/gin-gonic/gin"
)

const (

    // BestCompression provides the highest-level of compression
    BestCompression = gzip.BestCompression

    // BestSpeed focuses on speed rather than space saving
    BestSpeed = gzip.BestSpeed

    // DefaultCompression is the middle ground between speed and space
    DefaultCompression = gzip.DefaultCompression

    // NoCompression leaves the data as is
    NoCompression = gzip.NoCompression
)

// Gzip compresses the response data
func Gzip(level int) gin.HandlerFunc {
    return func(c *gin.Context) {

        // Create gzip writer from gin.ResponseWriter
        gz, err := gzip.NewWriterLevel(c.Writer, level)
        if err != nil {
            return
        }

        // Set response headers
        c.Header("Content-Encoding", "gzip")
        c.Header("Vary", "Accept-Encoding")

        // Overload response writer
        c.Writer = &gzipWriter{c.Writer, gz}
        defer func() {
            c.Header("Content-Length", "")
            gz.Close()
        }()

        // Call next middleware
        c.Next()
    }
}

type gzipWriter struct {
    gin.ResponseWriter
    writer *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
    return g.writer.Write(data)
}
