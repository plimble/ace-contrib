package gzip

import (
	"compress/gzip"
	"github.com/plimble/copter"
	"strings"
)

const (
	encodingGzip = "gzip"

	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
	headerContentLength   = "Content-Length"
	headerContentType     = "Content-Type"
	headerVary            = "Vary"

	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

type gzipWriter struct {
	copter.ResponseWriter
	gzwriter *gzip.Writer
}

func newGzipWriter(writer copter.ResponseWriter, gzwriter *gzip.Writer) *gzipWriter {
	return &gzipWriter{writer, gzwriter}
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.gzwriter.Write(data)
}

func Gzip(level int) copter.HandlerFunc {

	return func(c *copter.C) {
		req := c.Request
		if !strings.Contains(req.Header.Get(headerAcceptEncoding), encodingGzip) {
			c.Next()
			return
		}

		writer := c.Writer
		gz, err := gzip.NewWriterLevel(writer, level)
		if err != nil {
			c.Next()
			return
		}
		defer gz.Close()

		headers := writer.Header()
		headers.Set(headerContentEncoding, encodingGzip)
		headers.Set(headerVary, headerAcceptEncoding)

		gzwriter := newGzipWriter(c.Writer, gz)
		c.Writer = gzwriter
		c.Next()
		writer.Header().Del(headerContentLength)
	}
}
