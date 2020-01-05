package handlers

import (
	"context"

	"github.com/caledfwlch1/enlabtest/types"
)

func Game(ctx context.Context, data *types.BodyData) string {
	return "OK"
}
