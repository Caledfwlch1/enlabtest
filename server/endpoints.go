package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caledfwlch1/enlabtest/types"
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

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cancel()
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go serv.scheduler(ctx)

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

func validateRequest(rw http.ResponseWriter, request *http.Request) bool {
	if request.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}

	if request.Body == http.NoBody {
		rw.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func validateHeaders(rw http.ResponseWriter, request *http.Request, userId *uuid.UUID, srcType *string) bool {
	user := request.Header["User-Id"]
	if len(user) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(rw, "empty User-Id")
		return false
	}

	var err error
	*userId, err = uuid.Parse(user[0])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(rw, err)
		return false
	}

	srcTypes := request.Header["Source-Type"]
	if len(srcTypes) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(rw, "empty Source-Type")
		return false
	}

	if srcTypes[0] != "game" &&
		srcTypes[0] != "server" &&
		srcTypes[0] != "payment" {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(rw, "bad Source-Type")
		return false
	}

	*srcType = srcTypes[0]
	return true
}

func jsonRequest(rw http.ResponseWriter, request *http.Request, data *types.DataOperation, userId uuid.UUID) bool {
	data, err := types.ParseBody(request)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(rw, err)
		return false
	}
	_ = request.Body.Close()

	data.UserId = userId
	return true
}
