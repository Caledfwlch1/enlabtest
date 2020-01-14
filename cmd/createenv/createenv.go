package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	"github.com/caledfwlch1/enlabtest/types"

	"github.com/google/uuid"

	"github.com/caledfwlch1/enlabtest/db/postgres"
)

var (
	connStr = flag.String("c", "", "database connection string")
)

func main() {
	flag.Parse()

	dbs, err := sql.Open("postgres", *connStr)
	if err != nil {
		log.Fatalf("error open database %s", err)
	}

	if err = createSchema(dbs); err != nil {
		log.Fatalf("error creating schema %s", err)
	}

	db, err := postgres.NewDatabase(*connStr)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	userID, _ := uuid.Parse(types.TestUser)
	_, err = db.CreateUser(ctx, userID)
}

func createSchema(dbs *sql.DB) error {
	for _, query := range querys {
		_, err := dbs.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	querys = []string{
		`CREATE TABLE "user"
(
    user_id uuid NOT NULL,
    balance real NOT NULL DEFAULT 0,
    CONSTRAINT user_pkey PRIMARY KEY (user_id)
)
    WITH (
        OIDS=FALSE
    );
ALTER TABLE "user"
    OWNER TO docker;`,
		`CREATE TABLE transaction
(
    transaction_id uuid NOT NULL,
    state integer NOT NULL DEFAULT 0,
    amount real NOT NULL DEFAULT 0,
    user_id uuid NOT NULL,
    "timestamp" time with time zone NOT NULL DEFAULT now(),
    CONSTRAINT transaction_pkey PRIMARY KEY (transaction_id),
    CONSTRAINT transaction_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES "user" (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION ON DELETE NO ACTION
)
    WITH (
        OIDS=TRUE
    );
ALTER TABLE transaction
    OWNER TO docker;

CREATE INDEX timestamp_idx
    ON transaction
        USING btree
        ("timestamp");`,
	}
)
