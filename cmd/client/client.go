package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/caledfwlch1/enlabtest/client"

	"github.com/caledfwlch1/enlabtest/types"
)

var (
	connStr       = flag.String("c", "http://127.0.0.1:8080/request", "connection string")
	state         = flag.String("s", "win", "state")
	userID        = flag.String("u", types.TestUser, "user ID")
	srcType       = flag.String("r", "game", "source type")
	amount        = flag.Float64("a", 0, "amount")
	transactionID = flag.String("t", "", "transaction ID")
)

func main() {
	flag.Parse()

	httpClient := http.DefaultClient

	cln := client.NewClient(httpClient, *connStr, *srcType, *userID)

	if *transactionID == "" {
		*transactionID = uuid.New().String()
	}

	log.Printf("user: %s, srcType: %s, transaction: %s, state: %s\n", *userID, *srcType, *transactionID, *state)
	data := types.UserRequest{
		State:         *state,
		Amount:        strconv.FormatFloat(float64(*amount), 'f', 2, 32),
		TransactionId: *transactionID,
	}

	bal, err := cln.RequestToServer(&data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("user: %s, balance: %.2f", *userID, bal)
}
