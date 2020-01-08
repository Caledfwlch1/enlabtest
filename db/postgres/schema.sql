CREATE TABLE "user"
(
    user_id uuid NOT NULL,
    balance real NOT NULL DEFAULT 0,
    CONSTRAINT user_pkey PRIMARY KEY (user_id)
)
    WITH (
        OIDS=FALSE
    );
ALTER TABLE "user"
    OWNER TO docker;

CREATE TABLE transaction
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
        ("timestamp");

