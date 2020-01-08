package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/caledfwlch1/enlabtest/db"

	"github.com/caledfwlch1/enlabtest/handlers"
	"github.com/caledfwlch1/enlabtest/types"
)

type Config struct {
	Ip      string
	Port    string
	ConnStr string
}

func (s *Config) FullAddress() string {
	return s.Ip + ":" + s.Port
}

func (s Config) String() string {
	return s.Ip + ":" + s.Port
}

type server struct {
	db   db.Database
	stop func()
}

func (s *server) requestHandler(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if !validateRequest(rw, request) {
		return
	}

	var (
		userId  uuid.UUID
		srcType string
	)

	if !validateHeaders(rw, request, &userId, &srcType) {
		return
	}

	var data types.Transaction
	if !jsonRequest(rw, request, &data, userId) {
		return
	}

	ctx := request.Context()
	var (
		balance float32
		err     error
	)

	switch srcType {
	case "game":
		balance, err = handlers.Game(ctx, s.db, &data)
	case "server":
		balance, err = handlers.Server(ctx, s.db, &data)
	case "payment":
		balance, err = handlers.Payment(ctx, s.db, &data)
	default:
		err = types.ErrorUnknownSourceType
	}

	var out []byte
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		strErr := err.Error()
		out, _ = json.Marshal(strErr)
	} else {
		out, _ = json.Marshal(balance)
	}
	_, _ = fmt.Fprintln(rw, string(out))
}

func (s *server) scheduler(ctx context.Context) {
	ticker := time.NewTicker(postProcessingRollBackTimeout)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			handlers.RollBack(ctx, s.db)
		}
	}
}
