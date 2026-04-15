package middleware

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// etagStore maps request URI → last known ETag (server-wide, thread-safe).
var etagStore sync.Map

// etagWriter fully buffers the response body and captures the status code.
// Nothing is written to the wire until flush() is called, which allows us to
// set the ETag response header BEFORE any bytes are sent.
type etagWriter struct {
	gin.ResponseWriter
	buf    bytes.Buffer
	status int
}

func (w *etagWriter) WriteHeader(code int)  { w.status = code }
func (w *etagWriter) WriteHeaderNow()       {} // delay until flush
func (w *etagWriter) Write(b []byte) (int, error) { return w.buf.Write(b) }
func (w *etagWriter) WriteString(s string) (int, error) { return w.buf.WriteString(s) }
func (w *etagWriter) Written() bool         { return w.buf.Len() > 0 }
func (w *etagWriter) Size() int             { return w.buf.Len() }
func (w *etagWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

// flush commits the buffered response to the real ResponseWriter.
func (w *etagWriter) flush(etag string, send304 bool) {
	rw := w.ResponseWriter
	if send304 {
		rw.Header().Set("ETag", etag)
		rw.Header().Set("Cache-Control", "no-cache")
		rw.WriteHeader(http.StatusNotModified)
		rw.WriteHeaderNow()
		return
	}
	if etag != "" {
		rw.Header().Set("ETag", etag)
		rw.Header().Set("Cache-Control", "no-cache")
	}
	rw.WriteHeader(w.Status())
	if w.buf.Len() > 0 {
		rw.Write(w.buf.Bytes()) //nolint:errcheck
	} else {
		rw.WriteHeaderNow()
	}
}

// ETagMiddleware implements ETag/If-None-Match HTTP caching for GET requests.
//
// Flow:
//  1. If client sends If-None-Match matching stored ETag → respond 304 immediately.
//  2. Otherwise buffer the full response, compute SHA-1 ETag, store and respond 200+ETag.
func ETagMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		key := c.Request.RequestURI

		// Step 1: check If-None-Match
		if clientETag := c.GetHeader("If-None-Match"); clientETag != "" {
			if stored, ok := etagStore.Load(key); ok && stored.(string) == clientETag {
				c.AbortWithStatus(http.StatusNotModified)
				c.Header("ETag", clientETag)
				c.Header("Cache-Control", "no-cache")
				return
			}
		}

		// Step 2: buffer response
		ew := &etagWriter{ResponseWriter: c.Writer, status: http.StatusOK}
		c.Writer = ew
		c.Next()

		// Step 3: compute ETag and flush
		status := ew.Status()
		if status >= http.StatusOK && status < http.StatusMultipleChoices && ew.buf.Len() > 0 {
			hash := sha1.Sum(ew.buf.Bytes())
			etag := fmt.Sprintf(`"%x"`, hash)
			etagStore.Store(key, etag)
			ew.flush(etag, false)
		} else {
			ew.flush("", false)
		}
	}
}
