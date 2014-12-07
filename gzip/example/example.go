package main

import (
	"fmt"
	"github.com/plimble/copter"
	"github.com/plimble/copter-contrib/gzip"
	"time"
)

func main() {
	r := copter.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.GET("/ping", func(c *copter.C) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
