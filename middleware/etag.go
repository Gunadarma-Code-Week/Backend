package middleware

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
// This middleware fully buffers the response body to compute a strong ETag (SHA-1).
func ETagMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Buffer the response
		ew := &etagWriter{ResponseWriter: c.Writer, status: http.StatusOK}
		originalWriter := c.Writer
		c.Writer = ew

		// Step 1: Run the handler to get the actual data from DB/service
		c.Next()

		// Step 2: After handler finished, compute ETag of the actual response body
		status := ew.Status()
		// Only cache successful 200-OKish responses with a body
		if status >= http.StatusOK && status < http.StatusMultipleChoices && ew.buf.Len() > 0 {
			hash := sha1.Sum(ew.buf.Bytes())
			etag := fmt.Sprintf(`"%x"`, hash)

			// Step 3: Compare with client's If-None-Match
			clientETag := c.GetHeader("If-None-Match")
			if clientETag == etag {
				// Match! We can avoid sending the body.
				ew.flush(etag, true)
			} else {
				// No match or no ETag sent, send the new data + new ETag
				ew.flush(etag, false)
			}
		} else {
			// Response was an error or empty, just flush normally without ETag
			c.Writer = originalWriter
			ew.flush("", false)
		}
	}
}
