package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/caledfwlch1/enlabtest/db"
	"github.com/google/uuid"

	"github.com/caledfwlch1/enlabtest/types"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	connStr = "postgres://docker:docker@127.0.0.1/test?sslmode=disable"
)

func Test_postgres_ApplyTransaction(t *testing.T) {
	type args struct {
		d *types.Transaction
	}

	p, err := NewDatabase(connStr)
	require.NoError(t, err, "error open database")
	ctx := context.Background()

	userId, err := p.CreateUser(ctx)
	require.NoError(t, err, "error creating user")

	tests := []struct {
		name    string
		args    args
		wants   float32
		wantErr bool
	}{
		{
			name:  "win200",
			args:  args{d: types.NewDataOperation(userId, types.Win, 200)},
			wants: 200,
		},
		{
			name:  "lost-100",
			args:  args{d: types.NewDataOperation(userId, types.Lost, 100)},
			wants: 100,
		},
		{
			name:  "lost-100",
			args:  args{d: types.NewDataOperation(userId, types.Lost, 100)},
			wants: 0,
		},
		{
			name:    "lost-100",
			args:    args{d: types.NewDataOperation(userId, types.Lost, 100)},
			wants:   0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if _, err := p.ApplyTransaction(ctx, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("postgres.DoOperation() error = %v, wantErr %v", err, tt.wantErr)
			}

			balance, err := p.GetBalance(ctx, tt.args.d.UserID)
			assert.NoError(t, err, "error getting balance")
			assert.Equal(t, tt.wants, balance, "balance mismatch")

		})
	}
}

const numbRecs = 10

func Test_postgres_RollBackLastN(t *testing.T) {
	type args struct {
		task *types.RollBackTask
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "rollback",
			args: args{task: &types.RollBackTask{
				RecNumb: numbRecs,
				Odd:     true,
			}},
		},
	}

	dbs, err := NewDatabase(connStr)
	require.NoError(t, err, "error open database")
	ctx := context.Background()

	userId, err := dbs.CreateUser(ctx)
	require.NoError(t, err, "error creating user")

	createTransactions(t, ctx, dbs, userId)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rbd := prepareTestData(t, ctx, dbs, tt.args.task.Odd)

			if err := dbs.RollBackLastN(ctx, tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("postgres.RollBackLastN() error = %v, wantErr %v", err, tt.wantErr)
			}
			checkBalances(t, ctx, dbs, rbd)
		})
	}
}

type rollBackTestData struct {
	userBalance  map[uuid.UUID]balance
	transactions []*types.Transaction
}

type balance struct {
	real  float32
	delta float32
}

func prepareTestData(t *testing.T, ctx context.Context, db db.Database, odd bool) *rollBackTestData {
	dops, err := db.GetLastRecords(ctx, numbRecs*2)
	require.NoError(t, err, "error getting last records")
	dops = selectRecords(dops, odd)

	out := rollBackTestData{
		userBalance:  make(map[uuid.UUID]balance),
		transactions: dops,
	}

	for _, dop := range dops {
		b := out.userBalance[dop.UserID]
		b.delta += dop.GetAmount()
		out.userBalance[dop.UserID] = b
	}

	for u := range out.userBalance {
		bal, err := db.GetBalance(ctx, u)
		if err != nil {
			assert.NoError(t, err, "error getting balance")
		}
		b := out.userBalance[u]
		b.real = bal
		out.userBalance[u] = b
	}

	return &out
}

func checkBalances(t *testing.T, ctx context.Context, db db.Database, rbd *rollBackTestData) {
	for u, b := range rbd.userBalance {
		bal, err := db.GetBalance(ctx, u)
		if err != nil {
			assert.NoError(t, err, "error getting balance")
		}
		assert.Equal(t, b.real-b.delta, bal, fmt.Sprintf("user %s balance mismatch", u))
	}
}

func createTransactions(t *testing.T, ctx context.Context, db db.Database, userId uuid.UUID) {
	// we must be sure that the rollback of the balance will not lead to a negative balance
	numbTrans := numbRecs * 3
	amount := numbTrans*numbTrans/2 + numbTrans

	dop := types.Transaction{
		UserID: userId,
		State:  types.Win,
		Amount: float32(amount),
		ID:     uuid.New(),
	}
	_, err := db.ApplyTransaction(ctx, &dop)
	require.NoError(t, err)

	for i := 1; i <= numbTrans; i++ {
		dop.State = types.OperationState(i%2) + 1
		dop.Amount = float32(i)
		dop.Amount = dop.GetAmount()
		dop.ID = uuid.New()
		_, err = db.ApplyTransaction(ctx, &dop)
		require.NoError(t, err)
	}
}
