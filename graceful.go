package graceful

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const DefaultShutdownTimeout = 5 * time.Second

// A Server defines parameters for gracefully running an HTTP server.
// The zero value for Server is a valid configuration.
type Server struct {
	*http.Server

	// ShutdownTimeout is the maximum duration for shutting down the server.
	// A zero or negative value means there will be no timeout.
	ShutdownTimeout time.Duration

	errChan chan error
}

// init server
func (s *Server) init() {
	s.errChan = make(chan error, 1)
}

func (s *Server) load(opts []Option) {
	for _, opt := range opts {
		opt.apply(s)
	}
}

// ListenAndServe listens on the TCP network address s.Addr and then
// calls s.Server.ListenAndServe to handle requests on incoming connections.
func (s *Server) ListenAndServe(opts ...Option) error {
	s.init()
	s.load(opts)
	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.errChan <- fmt.Errorf("[graceful] %w", err)
		}
	}()
	return s.waitForShutdown()
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and
// then calls ServeTLS to handle requests on incoming TLS connections.
// Accepted connections are configured to enable TCP keep-alives.
func (s *Server) ListenAndServeTLS(certFile, keyFile string, opts ...Option) error {
	s.init()
	s.load(opts)
	go func() {
		if err := s.Server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			s.errChan <- fmt.Errorf("[graceful] %w", err)
		}
	}()
	return s.waitForShutdown()
}

// waiting for shutdown or error occur.
func (s *Server) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-s.errChan:
		return err
	case <-quit:
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
// ShutdownTimeout defaults 5 seconds.
func ListenAndServe(addr string, handler http.Handler, opts ...Option) error {
	server := &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		ShutdownTimeout: DefaultShutdownTimeout,
	}
	return server.ListenAndServe(opts...)
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects HTTPS connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate, any intermediates, and the CA's certificate.
// ShutdownTimeout defaults 5 seconds.
func ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler, opts ...Option) error {
	server := &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		ShutdownTimeout: DefaultShutdownTimeout,
	}
	return server.ListenAndServeTLS(certFile, keyFile, opts...)
}
