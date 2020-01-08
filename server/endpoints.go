package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"

	"github.com/caledfwlch1/enlabtest/db/postgres"
)

func newServer(conf *Config, f func()) (*server, error) {
	db, err := postgres.NewDatabase(conf.ConnStr)
	if err != nil {
		return nil, err
	}

	return &server{
		db:   db,
		stop: f,
	}, nil
}

func ListenAndServe(conf *Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serv, err := newServer(conf, cancel)
	if err != nil {
		return err
	}

	http.HandleFunc("/request", serv.requestHandler)

	httpSrv := &http.Server{
		Addr: conf.FullAddress(),
	}

	stop := make(chan error, 1)
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stop <- err
		}
	}()

	go serv.scheduler(ctx)

	log.Println("Server started")

	return shutdownServer(httpSrv, cancel, stop)
}

func shutdownServer(srv *http.Server, cancel context.CancelFunc, stop <-chan error) error {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-stop:
		cancel()
		return fmt.Errorf("listen: %v\n", err)
	case <-quit:
	}

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

func validateRequest(rw http.ResponseWriter, request *http.Request) bool {
	if request.Method != http.MethodPost {
		jsonError(rw, fmt.Errorf("bad method"), http.StatusBadRequest)
		return false
	}

	if request.Body == http.NoBody {
		jsonError(rw, fmt.Errorf("body is empty"), http.StatusBadRequest)
		return false
	}
	return true
}

func validateHeaders(rw http.ResponseWriter, request *http.Request, userId *uuid.UUID, srcType *string) bool {
	user := request.Header.Get("User-Id")

	var err error
	*userId, err = uuid.Parse(user)
	if err != nil {
		jsonError(rw, err, http.StatusBadRequest)
		return false
	}

	srcTypes := request.Header.Get("Source-Type")

	if srcTypes != "game" &&
		srcTypes != "server" &&
		srcTypes != "payment" {

		jsonError(rw, fmt.Errorf("bad Source-Type"), http.StatusBadRequest)
		return false
	}

	*srcType = srcTypes
	return true
}

func jsonRequest(rw http.ResponseWriter, req *http.Request, data interface{}) bool {
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(data)
	if err != nil {
		jsonError(rw, fmt.Errorf("error parsing body data: %s", err), http.StatusBadRequest)
		return false
	}
	return true
}

func jsonError(rw http.ResponseWriter, err error, statusCode int) {
	rw.WriteHeader(statusCode)
	enc := json.NewEncoder(rw)
	_ = enc.Encode(struct {
		Err string
	}{Err: err.Error()})
}
