package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/caledfwlch1/enlabtest/types"
	"github.com/google/uuid"
)

const (
	connStr = "http://localhost:8080/request"
)

type client struct {
	http    *http.Client
	userID  string
	srcType string
}

func main() {

	userID := "419032e5-d2b4-4711-b83d-77e0aed0e832"
	srcType := "game"
	client := http.DefaultClient

	cln := newClient(client, srcType, userID)

	log.Printf("user: %s, srcType: %s\n", userID, srcType)
	begBalance := 1000

	balance, err := cln.init(begBalance)
	if err != nil {
		log.Fatalln("init: ", err)
	}

	log.Printf("balance: %f\n", balance)

	for i := 0; i < 20; i++ {
		balance, err = cln.updateBalance(balance, i)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("balance: %f\n", balance)
	}
}

func newClient(clnt *http.Client, srcType, userID string) *client {
	return &client{
		http:    clnt,
		userID:  userID,
		srcType: srcType,
	}
}

func (c *client) init(bal int) (float32, error) {
	return c.requestToServer(0, bal)
}

func (c *client) updateBalance(balance float32, i int) (float32, error) {
	return c.requestToServer(balance, i)
}

func (c *client) requestToServer(balance float32, i int) (float32, error) {
	req, delta, err := makeHttpPostRequest(c.srcType, c.userID, i)
	if err != nil {
		return 0, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}

	bal, err := processResponse(resp)
	if err != nil {
		return 0, err
	}

	if delta+balance != bal {
		return 0, fmt.Errorf("balance mismatch %f != %f", delta+balance, bal)
	}

	return bal, nil
}

func processResponse(resp *http.Response) (float32, error) {
	if resp == nil {
		return 0, fmt.Errorf("empty response")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read error: %s", err)
	}

	if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("status code: %d, error: %s", resp.StatusCode, b)
	}

	var (
		bal    float32
		errStr string
	)
	err = json.Unmarshal(b, &bal)
	if err != nil {
		err = json.Unmarshal(b, &errStr)
		if err != nil {
			return 0, fmt.Errorf("unmarshal error: %s", err)
		}
		return 0, fmt.Errorf(errStr)
	}

	return float32(bal), nil
}

func makeHttpPostRequest(srcType, userID string, i int) (*http.Request, float32, error) {
	rd, bal, err := generateData(i)
	if err != nil {
		return nil, 0, fmt.Errorf("data generation error: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, connStr, rd)
	if err != nil {
		return nil, 0, fmt.Errorf("request generation error: %s", err)
	}

	req.Header.Set("Source-Type", srcType)
	req.Header.Set("User-Id", userID)

	return req, bal, nil
}

func generateData(i int) (io.Reader, float32, error) {
	state := types.OperationState(i%2 + 1)

	data := types.UserRequest{
		State:         state.String(),
		Amount:        strconv.FormatFloat(float64(i), 'f', 2, 32),
		TransactionId: uuid.New().String(),
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return nil, 0, err
	}

	if state == types.Lost {
		return bytes.NewReader(b), float32(-i), nil
	}
	return bytes.NewReader(b), float32(i), nil
}
