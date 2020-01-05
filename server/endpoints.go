package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func newServer() *Server {
	return &Server{}
}

func Load(conf *Config) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	serv := newServer()
	http.HandleFunc("/request", serv.requestHandler)

	httpSrv := &http.Server{
		Addr: conf.FullAddress(),
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server started")

	return shutdownServer(httpSrv, cancel)
}

func shutdownServer(srv *http.Server, cancel context.CancelFunc) error {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server ...")
	cancel()

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), timeoutServerShutdown)
	defer cancelTimeout()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		return fmt.Errorf("server shutdown: %s", err)
	}

	log.Println("Server exiting")
	return nil
}
