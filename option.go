package graceful

import (
	"time"
)

// Option interface
type Option interface {
	// apply option to server
	apply(*Server)
}

// option wraps a function that modifies Server into an
// implementation of the Option interface.
type option struct {
	f func(*Server)
}

func (o *option) apply(s *Server) {
	o.f(s)
}

func newOption(f func(s *Server)) *option {
	return &option{f: f}
}

// WithShutdownTimeout set timeout for shutting down server
func WithShutdownTimeout(timeout time.Duration) Option {
	return newOption(func(s *Server) {
		s.ShutdownTimeout = timeout
	})
}

// WithShutdownFunc registers function(s) running while shutting down by calling s.RegisterOnShutdown.
func WithShutdownFunc(fs ...func()) Option {
	return newOption(func(s *Server) {
		for _, f := range fs {
			s.RegisterOnShutdown(f)
		}
	})
}
