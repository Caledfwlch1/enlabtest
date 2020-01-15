package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/caledfwlch1/enlabtest/client"
	"github.com/caledfwlch1/enlabtest/server"
	"gopkg.in/yaml.v2"

	"github.com/caledfwlch1/enlabtest/types"
)

const (
	host           = "localhost"
	port           = "8080"
	connStr        = "http://" + host + ":" + port + "/request"
	configFileName = "server.yml"
)

func TestIntegration(t *testing.T) {
	go runServer()
	// give time to load the server
	time.Sleep(time.Second)

	userID := types.TestUser
	srcType := "game"
	httpClient := http.DefaultClient

	cln := client.NewClient(httpClient, connStr, srcType, userID)

	log.Printf("user: %s, srcType: %s\n", userID, srcType)
	begBalance := 1000

	balance, err := initUserBalance(cln, begBalance)
	if err != nil {
		t.Fatal("init: ", err)
	}

	t.Logf("balance: %f", balance)

	for i := 0; i < 20; i++ {
		state := types.OperationState(i%2 + 1)
		balance, err = updateBalance(cln, balance, float32(i), state)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("balance: %f", balance)
	}
}

func initUserBalance(c *client.AppClient, bal int) (float32, error) {
	data := types.UserRequest{
		State:         "win",
		Amount:        strconv.FormatFloat(float64(bal), 'f', 2, 32),
		TransactionId: uuid.New().String(),
	}
	return c.RequestToServer(&data)
}

func updateBalance(c *client.AppClient, balance, delta float32, state types.OperationState) (float32, error) {
	data := types.UserRequest{
		State:         state.String(),
		Amount:        strconv.FormatFloat(float64(delta), 'f', 2, 32),
		TransactionId: uuid.New().String(),
	}

	bal, err := c.RequestToServer(&data)
	if err != nil {
		return 0, err
	}

	delta *= state.Factor()
	if delta+balance != bal {
		return 0, fmt.Errorf("balance mismatch %f != %f", delta+balance, bal)
	}
	return bal, nil
}

func runServer() {
	saveConfigToFile(&server.Config{
		Ip:      host,
		Port:    port,
		ConnStr: types.ServerConnStr,
	})

	conf, err := server.NewConfig(host, port, types.ServerConnStr)
	if err != nil {
		log.Fatalln(err)
	}

	err = server.ListenAndServe(conf)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func saveConfigToFile(conf *server.Config) {
	b, err := yaml.Marshal(conf)
	if err != nil {
		log.Println("configuration data encoding error")
		return
	}

	err = ioutil.WriteFile(configFileName, b, 0666)
	if err != nil {
		log.Println("error writing to configuration file")
	}
}
