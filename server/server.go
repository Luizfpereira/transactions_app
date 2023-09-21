package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Port   string
	Router *gin.Engine
}

type HandlerDetails struct {
	HttpMethod string
	Handler    gin.HandlerFunc
}

func NewServer(port string, router *gin.Engine) *Server {
	return &Server{
		Port:   port,
		Router: router,
	}
}

func (s *Server) Start() {
	srv := &http.Server{
		Addr:    s.Port,
		Handler: s.Router,
	}

	go func() {
		log.Println("Listening and serving on port: ", s.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// will block the app until receives a signal from the O.S.
	<-quit

	// cancel will release all resources associated with the context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	// gracefully stop accepting new requests and waits for the active ones to be handled
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}
