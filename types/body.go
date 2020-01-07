package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/caledfwlch1/enlabtest/tools"

	"github.com/google/uuid"
)

type UserRequest struct {
	State         string `json:"state"`
	Amount        string `json:"amount"`
	TransactionId string `json:"transactionId"`
}

// for more accurate money operations,
// there is an option to save the “Amount” field as an integer value
type DataOperation struct {
	UserId        uuid.UUID `json:"-"`
	State         OperationState
	Amount        float32
	TransactionId uuid.UUID
}

type OperationState int

func NewDataOperation(userId uuid.UUID, state OperationState, amount float32) *DataOperation {
	return &DataOperation{
		UserId:        userId,
		State:         state,
		Amount:        amount,
		TransactionId: uuid.New(),
	}
}

func (d *DataOperation) GetAmount() float32 {
	return tools.IIF(d.State == Lost, -d.Amount, d.Amount).(float32)
}

func (d *DataOperation) GetInvertAmount() float32 {
	return -d.GetAmount()
}

func (o OperationState) String() string {
	return StateToString[o]
}

func ParseBody(request *http.Request) (*DataOperation, error) {
	var data DataOperation

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error parsing body data: %s", err)
	}

	return &data, nil
}

func (d *DataOperation) UnmarshalJSON(b []byte) error {
	var data UserRequest

	err := json.NewDecoder(bytes.NewReader(b)).Decode(&data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data: %s", err)
	}

	id, err := uuid.Parse(data.TransactionId)
	if err != nil {
		return fmt.Errorf("error unmarshaling transactionId field: %s", err)
	}
	d.TransactionId = id

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

func (d DataOperation) String() string {
	return fmt.Sprintf("userid:%s, state:%s, amount:%.2f, transactionId:%s",
		d.UserId, d.State, d.Amount, d.TransactionId)
}
