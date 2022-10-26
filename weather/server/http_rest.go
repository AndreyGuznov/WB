package server

import (
	"net/http"
	"weather/tool/log"
)

// HTTPRestServer ...
type HTTPRestServer struct {
	address string
	wrapped *http.Server
}

// NewHTTPRestServer creates HTTP service
func NewHTTPRestServer(address string, handler http.Handler) *HTTPRestServer {
	httpsrv := http.Server{
		Addr:    address,
		Handler: handler,
	}

	return &HTTPRestServer{wrapped: &httpsrv, address: address}
}

// Serve HTTP service

func (server *HTTPRestServer) Serve() error {
	log.Info("Serving on " + server.address)
	err := server.wrapped.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Err("Server crashed", err)
	}

	return err
}
