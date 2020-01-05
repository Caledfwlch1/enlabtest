package handlers

import (
	"context"

	"github.com/caledfwlch1/enlabtest/types"
)

func Server(ctx context.Context, data *types.BodyData) string {
	return "OK"
}
