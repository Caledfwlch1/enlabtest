package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/caledfwlch1/enlabtest/types"
)

type AppClient struct {
	http    *http.Client
	userID  string
	srcType string
	connStr string
}

func NewClient(clnt *http.Client, connStr, srcType, userID string) *AppClient {
	return &AppClient{
		http:    clnt,
		userID:  userID,
		srcType: srcType,
		connStr: connStr,
	}
}

func (c *AppClient) RequestToServer(data *types.UserRequest) (float32, error) {
	req, err := c.makeHttpPostRequest(data)
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

	var bal struct{ Balance float32 }

	err = json.Unmarshal(b, &bal)
	if err != nil {
		return 0, err
	}

	return bal.Balance, nil
}

func (c *AppClient) makeHttpPostRequest(data *types.UserRequest) (*http.Request, error) {
	rd, err := makeReader(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.connStr, rd)
	if err != nil {
		return nil, fmt.Errorf("request generation error: %s", err)
	}

	req.Header.Set("Source-Type", c.srcType)
	req.Header.Set("User-Id", c.userID)

	return req, nil
}

func makeReader(data *types.UserRequest) (io.Reader, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
