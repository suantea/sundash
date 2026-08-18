package middleware

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
)

// gzipResponseWriter compresses compressible content types, passing others
// through untouched (images etc. stay as-is).
type gzipResponseWriter struct {
	gin.ResponseWriter
	gz     *gzip.Writer
	passed bool // true once we decided to pass through uncompressed
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	if !w.passed {
		w.passed = true
		if isCompressible(w.Header().Get("Content-Type")) {
			// Length changes under gzip; let the writer omit it.
			w.Header().Del("Content-Length")
			return w.gz.Write(data)
		}
		w.Header().Del("Content-Encoding")
		return w.ResponseWriter.Write(data)
	}
	return w.gz.Write(data)
}

func (w *gzipResponseWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (w *gzipResponseWriter) Flush() {
	if w.passed {
		w.ResponseWriter.Flush()
		return
	}
	w.gz.Flush()
	w.ResponseWriter.Flush()
}

func isCompressible(contentType string) bool {
	ct := strings.ToLower(contentType)
	return strings.HasPrefix(ct, "text/") ||
		strings.Contains(ct, "json") ||
		strings.Contains(ct, "javascript") ||
		strings.Contains(ct, "xml") ||
		strings.Contains(ct, "svg") ||
		strings.Contains(ct, "font")
}

// Gzip compresses text/JSON responses for clients that accept gzip.
func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}
		// Skip responses that already declare an encoding.
		if c.Writer.Header().Get("Content-Encoding") != "" {
			c.Next()
			return
		}

		gz := gzip.NewWriter(c.Writer)
		c.Writer = &gzipResponseWriter{ResponseWriter: c.Writer, gz: gz}
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		defer gz.Close()

		c.Next()
	}
}
