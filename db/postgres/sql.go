package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/caledfwlch1/enlabtest/types"

	"github.com/caledfwlch1/enlabtest/db"
	_ "github.com/lib/pq"
)

type postgres struct {
	db *sql.DB
}

func NewDatabase(connStr string) (db.Database, error) {
	dbs, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error open database %s", err)
	}

	out := postgres{db: dbs}

	err = out.makeStoredProc()
	if err != nil {
		_ = dbs.Close()
		return nil, fmt.Errorf("error creating storage procedure: %s", err)
	}

	return &out, nil
}

func (p *postgres) ApplyTransaction(ctx context.Context, d *types.Transaction) (float32, error) {
	query := "SELECT * FROM update_user_balance($1, $2, $3, $4);"
	row := p.db.QueryRowContext(ctx, query, d.ID, d.State, d.GetAmount(), d.UserID)

	var result float32
	err := row.Scan(&result)
	if err != nil {
		return -1, err
	}
	if result < 0 {
		return result, fmt.Errorf("user balance cannot be negative")
	}

	return result, nil
}

func (p *postgres) makeStoredProc() error {
	for _, query := range storedProc {
		_, err := p.db.Exec(query)
		return err
	}
	return nil
}

func (p *postgres) GetBalance(ctx context.Context, userId uuid.UUID) (float32, error) {
	query := `SELECT balance FROM "user" WHERE user_id = $1;`

	row := p.db.QueryRowContext(ctx, query, userId)

	var balance float32
	err := row.Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (p *postgres) CreateUser(ctx context.Context) (uuid.UUID, error) {
	userId := uuid.New()

	query := `INSERT INTO "user" (user_id) VALUES ($1);`

	res, err := p.db.ExecContext(ctx, query, userId)
	if err != nil {
		return uuid.UUID{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return uuid.UUID{}, err
	}
	if rows != 1 {
		return uuid.UUID{}, fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}
	return userId, nil
}

func (p *postgres) RollBackLastN(ctx context.Context, task *types.RollBackTask) error {
	dops, err := p.GetLastRecords(ctx, task.RecNumb*2)
	if err != nil {
		return err
	}

	dops = selectRecords(dops, task.Odd)

	for _, dop := range dops {
		err = p.RollBackTransaction(ctx, &dop)
		if err != nil {
			log.Printf("error transaction roll back, transaction Id: %s, error: %s",
				dop.ID, err)
			continue
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
	return nil
}

func (p *postgres) GetLastRecords(ctx context.Context, n int) ([]types.Transaction, error) {
	query := `SELECT transaction_id, state, amount, user_id FROM transaction 
ORDER BY timestamp DESC LIMIT $1;`

	rows, err := p.db.QueryContext(ctx, query, n)
	if err != nil {
		return nil, err
	}

	var out []types.Transaction

	for rows.Next() {
		var (
			transactionId uuid.UUID
			state         int
			amount        float32
			userId        uuid.UUID
		)

		if err := rows.Scan(&transactionId, &state, &amount, &userId); err != nil {
			continue
		}

		out = append(out, types.Transaction{
			ID:     transactionId,
			State:  types.OperationState(state),
			Amount: amount,
			UserID: userId,
		})
	}

	return out, nil
}

func selectRecords(dops []types.Transaction, odd bool) []types.Transaction {
	var (
		out   []types.Transaction
		start int
	)

	if !odd {
		start = 1
	}

	for i := start; i < len(dops); i += 2 {
		out = append(out, dops[i])
	}
	return out
}

func (p *postgres) RollBackTransaction(ctx context.Context, dop *types.Transaction) error {
	query := `select * from rollback_transaction($1, $2, $3);`
	row := p.db.QueryRowContext(ctx, query, dop.ID, -dop.GetAmount(), dop.UserID)

	var bal float32
	err := row.Scan(&bal)
	if err != nil {
		return err
	}

	if bal < 0 {
		return fmt.Errorf("rollback is not possible - the balance cannot be negative")
	}

	return nil
}
