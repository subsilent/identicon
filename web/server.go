package web

import (
	"encoding/base64"
	"time"

	identicon "github.com/dgryski/go-identicon"
	"github.com/gin-gonic/gin"
	log "github.com/mgutz/logxi/v1"
)

var Key = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

func NewServer(logger log.Logger) *gin.Engine {

	// create router
	router := gin.New()
	router.Use(gin.Recovery(), Logger(logger))
	router.GET("/icon/:name", BasicImageRoute(logger))
	router.GET("/encoded/:name", EncodedImageRoute(logger))
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

func EncodedImageRoute(l log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		icon := identicon.New7x7(Key)
		name := c.Param("name")

		l.Info("Encoding identicon", "name", name)

		data := []byte(name)
		pngdata := icon.Render(data)
		encoding := base64.StdEncoding.EncodeToString(pngdata)
		c.Writer.Write([]byte("data:image/png;base64," + encoding))
		c.Header("Content-Type", "text/plain")
	}
}
