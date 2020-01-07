package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/caledfwlch1/enlabtest/db"

	"github.com/caledfwlch1/enlabtest/handlers"
	"github.com/caledfwlch1/enlabtest/types"
)

type Config struct {
	Ip       string
	Port     string
	Host     string
	User     string
	Pass     string
	Database string
	Options  string
}

func (s *Config) FullAddress() string {
	return s.Ip + ":" + s.Port
}

func (s Config) String() string {
	return s.Ip + ":" + s.Port
}

type Server struct {
	db   db.Database
	stop chan struct{}
}

func (s *Server) requestHandler(rw http.ResponseWriter, request *http.Request) {
	if err := validateRequest(request); err != nil {
		_, _ = fmt.Fprintln(rw, err)
		return
	}

	ctx := request.Context()

	data, err := types.ParseBody(request)
	if err != nil {
		_, _ = fmt.Fprintln(rw, err)
		return
	}
	_ = request.Body.Close()

	user := request.Header["User-Id"]
	if len(user) == 0 {
		_, _ = fmt.Fprintln(rw, "empty User-Id")
		return
	}
	userId, err := uuid.Parse(user[0])
	if err != nil {
		_, _ = fmt.Fprintln(rw, err)
		return
	}
	data.UserId = userId

	srcType := request.Header["Source-Type"]
	if len(srcType) == 0 {
		_, _ = fmt.Fprintln(rw, "empty Source-Type")
		return
	}

	var resp string

	switch srcType[0] {
	case "game":
		resp = handlers.Game(ctx, s.db, data)
	case "server":
		resp = handlers.Server(ctx, s.db, data)
	case "payment":
		resp = handlers.Payment(ctx, s.db, data)
	default:
		resp = "unknown Source-Type"
	}

	_, _ = fmt.Fprintln(rw, resp)
}

func (s *Server) Scheduler() {
	ctx := context.Background()

	go func() {
		<-s.stop
		ctx.Done()
	}()

	// RollBack task
	go func() {
		ticker := time.NewTicker(postProcessingRollBackTimeout)
		defer ticker.Stop()
		for {
			select {
			case <-s.stop:
				return

			case <-ticker.C:
				handlers.RollBack(ctx, s.db)
			}
		}
	}()

	// we can add more tasks here
	<-s.stop
}
