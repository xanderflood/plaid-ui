package db

import (
	"context"

	"github.com/pkg/errors"
)

func (a *DBAgent) EnsureTransactionsTable(ctx context.Context) error {
	_, err := a.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS "transactions"
(	"uuid" UUID,
	"account_uuid" UUID REFERENCES accounts(uuid),
	"user_uuid" UUID REFERENCES users(uuid),
	"created_at" timestamp NOT NULL,
	"modified_at" timestamp NOT NULL,
	"deleted_at" timestamp,

	"iso_currency_code" varchar,
	"amount" decimal,
	"date" varchar,

	"plaid_account_id" varchar,
	"plaid_name" varchar,
	"plaid_category_id" varchar,
	"plaid_pending" boolean,
	"plaid_pending_transaction_id" varchar,
	"plaid_account_owner" varchar,
	"plaid_transaction_id" varchar,
	"plaid_type" varchar,
	PRIMARY KEY ("uuid")
)`)
	if err != nil {
		return errors.Wrapf(err, "failed to ensure transactions table")
	}

	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON transactions USING btree(account_uuid)`)
	if err != nil {
		return errors.Wrap(err, "failed to ensure account_uuid index for transactions")
	}

	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON transactions USING btree(plaid_transaction_id)`)
	if err != nil {
		return errors.Wrap(err, "failed to ensure plaid_transaction_id index for transactions")
	}

	return nil
}

func (a *DBAgent) UpsertTransaction(ctx context.Context, accountUUID string, transaction Transaction) (bool, error) {
	uuid := a.uuider.UUID()
	row := a.db.QueryRowContext(ctx, `
INSERT INTO "transactions" (
	"uuid",
	"account_uuid",
	"user_uuid",
	"created_at",
	"modified_at",

	"iso_currency_code",
	"amount",
	"date",

	"plaid_account_id",
	"plaid_name",
	"plaid_category_id",
	"plaid_pending",
	"plaid_pending_transaction_id",
	"plaid_account_owner",
	"plaid_transaction_id",
	"plaid_type"
) VALUES (
	$1, $2, $3, NOW(), NOW(),
	$4, $5, $6,
	$7, $8, $9, $10, $11, $12, $13, $14
) ON CONFLICT ("uuid")
DO UPDATE SET
	"modified_at" = NOW(),
	"amount" = $5,
	"plaid_pending" = $10,
	"plaid_pending_transaction_id" = $11
RETURNING "created_at" = "modified_at"`,
		uuid,
		transaction.AccountUUID,
		transaction.UserUUID,

		transaction.ISOCurrencyCode,
		transaction.Amount,
		transaction.Date,

		transaction.PlaidAccountID,
		transaction.PlaidName,
		transaction.PlaidCategoryID,
		transaction.PlaidPending,
		transaction.PlaidPendingTransactionID,
		transaction.PlaidAccountOwner,
		transaction.PlaidID,
		transaction.PlaidType,
	)

	var isNew bool
	err := row.Scan(&isNew)
	return isNew, errors.Wrapf(err, "failed to upsert to transactions table for plaid transaction %s", transaction.PlaidID)
}

func (a *DBAgent) DeleteTransactionByPlaidID(ctx context.Context, plaidTransactionID string) error {
	_, err := a.db.ExecContext(ctx, `
UPDATE "accounts"
SET "deleted_at" = NOW()
WHERE "plaid_transaction_id" = $1`,
		plaidTransactionID,
	)
	return errors.Wrapf(err, "failed to insert into transactions table")
}

func (a *DBAgent) GetTransactions(ctx context.Context, userUUID string, accountUUID string) ([]Transaction, error) {
	//TODO pagination
	rows, err := a.db.QueryContext(ctx, `
SELECT
	"uuid",
	"account_uuid",
	"user_uuid",
	"created_at",
	"modified_at",

	"iso_currency_code",
	"amount",
	"date",

	"plaid_account_id",
	"plaid_name",
	"plaid_category_id",
	"plaid_pending",
	"plaid_pending_transaction_id",
	"plaid_account_owner",
	"plaid_transaction_id",
	"plaid_type"
FROM "accounts"
WHERE
	"deleted_at" IS NULL
	AND "user_uuid" = $1
	AND "account_uuid" = $2
`,
		userUUID,
		accountUUID,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get transactions from table")
	}

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		err = rows.Scan(
			&transaction.UUID,
			&transaction.AccountUUID,
			&transaction.UserUUID,
			&transaction.CreatedAt,
			&transaction.ModifiedAt,

			&transaction.ISOCurrencyCode,
			&transaction.Amount,
			&transaction.Date,

			&transaction.PlaidAccountID,
			&transaction.PlaidName,
			&transaction.PlaidCategoryID,
			&transaction.PlaidPending,
			&transaction.PlaidPendingTransactionID,
			&transaction.PlaidAccountOwner,
			&transaction.PlaidID,
			&transaction.PlaidType,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to scan result of querying for all transactions for account %s", accountUUID)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
