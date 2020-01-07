package handlers

import (
	"context"

	"github.com/caledfwlch1/enlabtest/db"

	"github.com/caledfwlch1/enlabtest/types"
)

func Game(ctx context.Context, db db.Database, data *types.DataOperation) string {
	if data.State != types.Win && data.State != types.Lost {
		return types.UnknownState
	}

	err := db.DoOperation(ctx, data)
	if err != nil {
		return err.Error()
	}

	return types.OperationOk
}
