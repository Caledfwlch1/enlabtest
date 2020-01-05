package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BodyData struct {
	State         string
	Amount        string
	TransactionId string
}

func NewBody(request *http.Request) (*BodyData, error) {
	var data BodyData

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error parsing body data: %s", err)
	}

	return &data, nil
}
