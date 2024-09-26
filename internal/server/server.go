package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (srv *Server) Run(port string, handler http.Handler) error {
	srv.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, //1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return srv.httpServer.ListenAndServe()
}

func (srv *Server) ShutDown(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}
