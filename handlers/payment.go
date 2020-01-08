package handlers

import (
	"context"

	"github.com/caledfwlch1/enlabtest/db"

	"github.com/caledfwlch1/enlabtest/types"
)

func Payment(ctx context.Context, db db.Database, data *types.Transaction) (float32, error) {
	if data.State != types.Win {
		return -1, types.ErrorUnknownState
	}
	return db.ApplyTransaction(ctx, data)
}
