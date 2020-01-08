package handlers

import (
	"context"
	"log"

	"github.com/caledfwlch1/enlabtest/db"
	"github.com/caledfwlch1/enlabtest/types"
)

func RollBack(ctx context.Context, db db.Database) {
	task := &types.RollBackTask{
		RecNumb: 10,
		Odd:     true,
	}
	err := db.RollBackLastN(ctx, task)
	if err != nil {
		log.Println("roll back error: ", err)
	}
}
