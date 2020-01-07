package postgres

var storedProc = []string{
	// this is a request to create a stored procedure to update user balance
	`CREATE OR REPLACE FUNCTION update_user_balance(trans UUID, state INTEGER, amount REAL, id_user UUID) RETURNS REAL AS $$
DECLARE 
   balance_before real;
   balance_after real;
BEGIN
	BEGIN
		SELECT INTO balance_before balance FROM "user" WHERE user_id = id_user;
		INSERT INTO "transaction" (transaction_id, state, amount, user_id) VALUES (trans, state, amount, id_user);
		UPDATE "user" SET balance = balance + amount WHERE user_id = id_user;
		SELECT INTO balance_after balance FROM "user" WHERE user_id = id_user;

		IF balance_before + amount != balance_after OR balance_after < 0 THEN
			RAISE 'check_violation';
		END IF;

	EXCEPTION WHEN OTHERS THEN
		RETURN -1;
	END;
RETURN balance_after;
END;
$$ LANGUAGE plpgsql;`,

	// this is a request to roll back a transaction and update the user balance
	`CREATE OR REPLACE FUNCTION rollback_transaction(trans UUID, amount REAL, id_user UUID) RETURNS REAL AS $$
DECLARE 
   balance_before real;
   balance_after real;
BEGIN
	BEGIN
		SELECT INTO balance_before balance FROM "user" WHERE user_id = id_user;
		DELETE FROM "transaction" WHERE transaction_id = trans;
		UPDATE "user" SET balance = balance + amount WHERE user_id = id_user;
		SELECT INTO balance_after balance FROM "user" WHERE user_id = id_user;

		IF balance_before + amount != balance_after OR balance_after < 0 THEN
			RAISE 'check_violation';
		END IF;

	EXCEPTION WHEN OTHERS THEN
		RETURN -1;
	END;
RETURN balance_after;
END;
$$ LANGUAGE plpgsql;`,
}
