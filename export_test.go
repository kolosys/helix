package helix

// Export unexported symbols for testing.

// NewRouter exports newRouter for testing.
func NewRouter() *Router {
	return newRouter()
}

// ParsePattern exports parsePattern for testing.
// Returns an int (the segment count) since segment is unexported.
func ParsePattern(pattern string) []any {
	segs := parsePattern(pattern)
	result := make([]any, len(segs))
	for i, s := range segs {
		result[i] = s
	}
	return result
}

// IsBindingError exports isBindingError for testing.
func IsBindingError(err error) bool {
	return isBindingError(err)
}

// ServerConfig holds server configuration for testing.
type ServerConfig struct {
	Addr           string
	ReadTimeout    int64
	WriteTimeout   int64
	IdleTimeout    int64
	GracePeriod    int64
	MiddlewareLen  int
	TLSCertFile    string
	TLSKeyFile     string
	MaxHeaderBytes int
}

// GetConfig returns the server configuration for testing.
func (s *Server) GetConfig() ServerConfig {
	return ServerConfig{
		Addr:           s.addr,
		ReadTimeout:    int64(s.readTimeout),
		WriteTimeout:   int64(s.writeTimeout),
		IdleTimeout:    int64(s.idleTimeout),
		GracePeriod:    int64(s.gracePeriod),
		MiddlewareLen:  len(s.middleware),
		TLSCertFile:    s.tlsCertFile,
		TLSKeyFile:     s.tlsKeyFile,
		MaxHeaderBytes: s.maxHeaderBytes,
	}
}
