package main

import (
	"fmt"
	"github.com/plimble/ace"
	"github.com/plimble/ace-contrib/gzip"
	"time"
)

func main() {
	a := ace.New()
	a.Use(gzip.Gzip(gzip.DefaultCompression))
	a.GET("/ping", func(c *ace.C) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	a.Run(":8080")
}
