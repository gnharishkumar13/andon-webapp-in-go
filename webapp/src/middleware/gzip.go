package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipMiddleware struct {
	wrapped    http.Handler
	exceptions []string
}

func NewGzip(exceptions []string, wrapped http.Handler) http.Handler {
	if wrapped == nil {
		wrapped = http.DefaultServeMux
	}
	if exceptions == nil {
		exceptions = make([]string, 0)
	}
	return &gzipMiddleware{
		wrapped:    wrapped,
		exceptions: exceptions,
	}
}

func (g gzipMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//use default when the exceptions ara available, applicable to websockets
	for _, e := range g.exceptions {
		if strings.HasPrefix(r.URL.Path, e) {
			g.wrapped.ServeHTTP(w, r)
			return
		}
	}

	//if browser does not support gzip
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		g.wrapped.ServeHTTP(w, r)
	}
	w.Header().Add("Content-Encoding", "gzip")

	//The below lines add the gzip to the response writer
	gzw := gzip.NewWriter(w)
	defer gzw.Close()

	grw := gzipResponseWriter{w, gzw}
	g.wrapped.ServeHTTP(grw, r)
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (grw gzipResponseWriter) Write(data []byte) (int, error) {
	return grw.Writer.Write(data)
}
