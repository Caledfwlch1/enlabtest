package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caledfwlch1/enlabtest/db/postgres"
)

func newServer(conf *Config) (*Server, error) {
	db, err := postgres.NewDatabase(conf.Host, conf.User, conf.Pass, conf.Database, conf.Options)
	if err != nil {
		return nil, err
	}

	return &Server{db: db}, nil
}

func Load(conf *Config) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	serv, err := newServer(conf)
	if err != nil {
		return err
	}

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

func validateRequest(request *http.Request) error {
	if request.Method != http.MethodPost {
		return fmt.Errorf("bad method")
	}

	if request.Body == http.NoBody {
		return fmt.Errorf("empty request")
	}
	return nil
}
