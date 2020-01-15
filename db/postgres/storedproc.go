package postgres

var storedProc = []string{
	// this is a request to create a stored procedure to update user balance
	`
CREATE OR REPLACE FUNCTION update_user_balance(trans UUID, state INTEGER, amount REAL, id_user UUID) RETURNS real AS $$
DECLARE 
   balance_before real;
   balance_after real;
   balance_temp real;
BEGIN
	BEGIN
		SELECT INTO balance_before balance FROM "user" WHERE user_id = id_user;
		balance_temp := balance_before + amount;

		IF balance_temp < 0 THEN
			RAISE EXCEPTION 'balance can not be negative %', balance_temp;
		END IF;

		INSERT INTO "transaction" (transaction_id, state, amount, user_id) VALUES (trans, state, amount, id_user);
		UPDATE "user" SET balance = balance_temp WHERE user_id = id_user;
		SELECT INTO balance_after balance FROM "user" WHERE user_id = id_user;

        IF balance_temp != balance_after THEN
            RAISE EXCEPTION 'balance missmatch % != %', balance_temp, balance_after;
        END IF;

		IF balance_after < 0 THEN
			RAISE EXCEPTION 'after update operation balance can not be negative %', balance_after;
		END IF;		
	END;

RETURN balance_after;
END;
$$ LANGUAGE plpgsql;
`,

	// this is a request to roll back a transaction and update the user balance
	`
CREATE OR REPLACE FUNCTION rollback_transaction(trans UUID, amount REAL, id_user UUID) RETURNS REAL AS $$
DECLARE 
   balance_before real;
   balance_after real;
   balance_temp real;
BEGIN
	BEGIN
		SELECT INTO balance_before balance FROM "user" WHERE user_id = id_user;
        balance_temp := balance_before - amount;

		IF balance_temp < 0 THEN
			RAISE EXCEPTION 'balance can not be negative %', balance_temp;
		END IF;

		DELETE FROM "transaction" WHERE transaction_id = trans;
		UPDATE "user" SET balance = balance_temp WHERE user_id = id_user;
		SELECT INTO balance_after balance FROM "user" WHERE user_id = id_user;

		IF balance_temp != balance_after THEN
            RAISE EXCEPTION 'balance missmatch % != %', balance_temp, balance_after;
        END IF;

		IF balance_after < 0 THEN
			RAISE EXCEPTION 'after update operation balance can not be negative %', balance_after;
		END IF;
	END;

RETURN balance_after;
END;
$$ LANGUAGE plpgsql;
`,
}

var dropStoredProc = []string{
	`DROP FUNCTION update_user_balance(uuid, integer, real, uuid);`,
	`DROP FUNCTION rollback_transaction(uuid, real, uuid);`,
}
