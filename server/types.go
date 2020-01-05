package server

import (
	"fmt"
	"net/http"

	"github.com/caledfwlch1/enlabtest/types"

	"github.com/caledfwlch1/enlabtest/handlers"
)

type Config struct {
	Ip   string
	Port string
}

func (s *Config) FullAddress() string {
	return s.Ip + ":" + s.Port
}

func (s Config) String() string {
	return s.Ip + ":" + s.Port
}

type Server struct{}

func (s *Server) requestHandler(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	srcType := request.Header["Source-Type"]
	if len(srcType) == 0 {
		_, _ = fmt.Fprintln(rw, "empty Source-Type")
		return
	}

	if request.Body == http.NoBody {
		_, _ = fmt.Fprintln(rw, "empty request")
		return
	}

	data, err := types.NewBody(request)
	if err != nil {
		_, _ = fmt.Fprintln(rw, err)
		return
	}
	_ = request.Body.Close()

	var resp string

	switch srcType[0] {
	case "game":
		resp = handlers.Game(ctx, data)
	case "server":
		resp = handlers.Server(ctx, data)
	case "payment":
		resp = handlers.Payment(ctx, data)
	default:
		resp = "unknown Source-Type"
	}

	_, _ = fmt.Fprintln(rw, resp)
}
