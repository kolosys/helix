package middleware

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

// CompressConfig configures the Compress middleware.
type CompressConfig struct {
	// Level is the compression level.
	// Valid levels: -1 (default), 0 (no compression), 1 (best speed) to 9 (best compression)
	// Default: -1 (gzip.DefaultCompression)
	Level int

	// MinSize is the minimum size in bytes to trigger compression.
	// Default: 1024 (1KB)
	MinSize int

	// Types is a list of content types to compress.
	// Default: text/*, application/json, application/javascript, application/xml
	Types []string

	// SkipFunc is a function that determines if compression should be skipped.
	SkipFunc func(r *http.Request) bool
}

// DefaultCompressConfig returns the default Compress configuration.
func DefaultCompressConfig() CompressConfig {
	return CompressConfig{
		Level:   gzip.DefaultCompression,
		MinSize: 1024,
		Types: []string{
			"text/",
			"application/json",
			"application/javascript",
			"application/xml",
			"application/xhtml+xml",
			"image/svg+xml",
		},
	}
}

// Compress returns a middleware that compresses responses using gzip or deflate.
func Compress() Middleware {
	return CompressWithConfig(DefaultCompressConfig())
}

// CompressWithLevel returns a Compress middleware with the given compression level.
func CompressWithLevel(level int) Middleware {
	config := DefaultCompressConfig()
	config.Level = level
	return CompressWithConfig(config)
}

// CompressWithConfig returns a Compress middleware with the given configuration.
func CompressWithConfig(config CompressConfig) Middleware {
	if config.Level < -1 || config.Level > 9 {
		config.Level = gzip.DefaultCompression
	}
	if config.MinSize <= 0 {
		config.MinSize = 1024
	}
	if len(config.Types) == 0 {
		config.Types = DefaultCompressConfig().Types
	}

	// Create pools for gzip and deflate writers
	gzipPool := &sync.Pool{
		New: func() any {
			w, _ := gzip.NewWriterLevel(io.Discard, config.Level)
			return w
		},
	}

	flatePool := &sync.Pool{
		New: func() any {
			w, _ := flate.NewWriter(io.Discard, config.Level)
			return w
		},
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if configured
			if config.SkipFunc != nil && config.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Check Accept-Encoding header
			acceptEncoding := r.Header.Get("Accept-Encoding")
			if acceptEncoding == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Determine encoding
			var encoding string
			if strings.Contains(acceptEncoding, "gzip") {
				encoding = "gzip"
			} else if strings.Contains(acceptEncoding, "deflate") {
				encoding = "deflate"
			} else {
				next.ServeHTTP(w, r)
				return
			}

			// Create compress writer
			cw := &compressWriter{
				ResponseWriter: w,
				encoding:       encoding,
				config:         config,
				gzipPool:       gzipPool,
				flatePool:      flatePool,
			}

			defer cw.Close()

			next.ServeHTTP(cw, r)
		})
	}
}

// compressWriter wraps http.ResponseWriter with compression.
type compressWriter struct {
	http.ResponseWriter
	encoding      string
	config        CompressConfig
	gzipPool      *sync.Pool
	flatePool     *sync.Pool
	writer        io.Writer
	gzipWriter    *gzip.Writer
	flateWriter   *flate.Writer
	buffer        []byte
	headerWritten bool
	compressed    bool
	statusCode    int
}

func (cw *compressWriter) WriteHeader(code int) {
	cw.statusCode = code
	// Don't write header yet - we need to check content type
}

func (cw *compressWriter) Write(b []byte) (int, error) {
	// Buffer until we have enough data
	cw.buffer = append(cw.buffer, b...)

	// Check if we should start compression
	if !cw.headerWritten && len(cw.buffer) >= cw.config.MinSize {
		cw.finalize()
	}

	return len(b), nil
}

func (cw *compressWriter) finalize() {
	if cw.headerWritten {
		return
	}
	cw.headerWritten = true

	contentType := cw.Header().Get("Content-Type")
	if contentType != "" && cw.shouldCompress(contentType) && len(cw.buffer) >= cw.config.MinSize {
		cw.startCompression()
	}

	if cw.statusCode == 0 {
		cw.statusCode = http.StatusOK
	}
	cw.ResponseWriter.WriteHeader(cw.statusCode)
}

func (cw *compressWriter) shouldCompress(contentType string) bool {
	for _, t := range cw.config.Types {
		if strings.HasSuffix(t, "/") {
			// Prefix match (e.g., "text/")
			if strings.HasPrefix(contentType, t) {
				return true
			}
		} else {
			// Starts with match
			if strings.HasPrefix(contentType, t) {
				return true
			}
		}
	}
	return false
}

func (cw *compressWriter) startCompression() {
	cw.compressed = true
	cw.Header().Set("Content-Encoding", cw.encoding)
	cw.Header().Del("Content-Length")
	cw.Header().Add("Vary", "Accept-Encoding")

	switch cw.encoding {
	case "gzip":
		gw := cw.gzipPool.Get().(*gzip.Writer)
		gw.Reset(cw.ResponseWriter)
		cw.gzipWriter = gw
		cw.writer = gw
	case "deflate":
		fw := cw.flatePool.Get().(*flate.Writer)
		fw.Reset(cw.ResponseWriter)
		cw.flateWriter = fw
		cw.writer = fw
	}
}

func (cw *compressWriter) Close() error {
	// Finalize if not yet done
	if !cw.headerWritten {
		cw.finalize()
	}

	// Write buffered data
	if len(cw.buffer) > 0 {
		if cw.compressed && cw.writer != nil {
			cw.writer.Write(cw.buffer)
		} else {
			cw.ResponseWriter.Write(cw.buffer)
		}
	}

	// Close compression writers and return to pool
	if cw.gzipWriter != nil {
		cw.gzipWriter.Close()
		cw.gzipPool.Put(cw.gzipWriter)
	}
	if cw.flateWriter != nil {
		cw.flateWriter.Close()
		cw.flatePool.Put(cw.flateWriter)
	}

	return nil
}

func (cw *compressWriter) Flush() {
	if cw.gzipWriter != nil {
		cw.gzipWriter.Flush()
	}
	if cw.flateWriter != nil {
		cw.flateWriter.Flush()
	}
	if f, ok := cw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
