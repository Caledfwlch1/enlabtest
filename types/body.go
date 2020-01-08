package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type UserRequest struct {
	State         string `json:"state"`
	Amount        string `json:"amount"`
	TransactionId string `json:"transactionId"`
}

// for more accurate money operations,
// there is an option to save the “Amount” field as an integer value
type Transaction struct {
	UserID uuid.UUID `json:"-"`
	State  OperationState
	Amount float32
	ID     uuid.UUID
}

type OperationState int

func NewDataOperation(userId uuid.UUID, state OperationState, amount float32) *Transaction {
	return &Transaction{
		UserID: userId,
		State:  state,
		Amount: amount,
		ID:     uuid.New(),
	}
}

func (d *Transaction) GetAmount() float32 {
	if d.State == Lost {
		return -d.Amount
	}
	return d.Amount
}

func (o OperationState) String() string {
	return StateToString[o]
}

func (d *Transaction) UnmarshalJSON(b []byte) error {
	var data UserRequest

	err := json.NewDecoder(bytes.NewReader(b)).Decode(&data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data: %s", err)
	}

	id, err := uuid.Parse(data.TransactionId)
	if err != nil {
		return fmt.Errorf("error unmarshaling transactionId field: %s", err)
	}
	d.ID = id

	switch data.State {
	case "win":
		d.State = Win
	case "lost":
		d.State = Lost
	default:
		return fmt.Errorf("error unmarshaling state field: %s", err)
	}

	amount, err := strconv.ParseFloat(data.Amount, 32)
	if err != nil {
		return fmt.Errorf("error parsing float value: %s", err)
	}
	d.Amount = float32(amount)
	return nil
}

func (d Transaction) String() string {
	return fmt.Sprintf("userid:%s, state:%s, amount:%.2f, transactionId:%s",
		d.UserID, d.State, d.Amount, d.ID)
}
