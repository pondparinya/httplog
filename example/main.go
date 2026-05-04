package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pondparinya/httplog"
	httploggin "github.com/pondparinya/httplog/gin"
)

func main() {
	l, err := httplog.NewZapLogger(httplog.ZapConfig{
		Level:               "debug",
		Format:              "json",
		DisableFunctionName: true,
	})
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(httploggin.Logger(l, []string{"/health"}))
	r.Use(httploggin.Recovery(l))

	r.GET("/hello", func(c *gin.Context) {
		ctx := c.Request.Context()
		l.Infocf(ctx, "handling hello request")
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	l.Infof("Server starting on :8080")
	r.Run(":8080")
}
