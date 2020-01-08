package db

import (
	"context"

	"github.com/google/uuid"

	"github.com/caledfwlch1/enlabtest/types"
)

type Database interface {
	CreateUser(ctx context.Context) (uuid.UUID, error)
	ApplyTransaction(ctx context.Context, d *types.Transaction) (float32, error)
	GetBalance(ctx context.Context, userId uuid.UUID) (float32, error)
	RollBackLastN(ctx context.Context, task *types.RollBackTask) error
	RollBackTransaction(ctx context.Context, td *types.Transaction) error
	GetLastRecords(ctx context.Context, n int) ([]types.Transaction, error)
}
