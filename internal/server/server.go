package server

import (
	"manga-library/pkg/logger"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

func NewServer(logger logger.Logger) *Server {
	return &Server{logger: logger}
}

const (
	maxHeaderBytes = 1 << 20
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
)

func (s *Server) Run(port, bindIP string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           bindIP + ":" + port,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}

	s.logger.Infof("server listening on: %s:%s", bindIP, port)
	return s.httpServer.ListenAndServe()
}
