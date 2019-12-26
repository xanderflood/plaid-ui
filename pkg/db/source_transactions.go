package db

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func (a *DBAgent) EnsureSourceTransactionsTable(ctx context.Context) error {
	_, err := a.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS "source_transactions"
(	"uuid" UUID DEFAULT gen_random_uuid(),
	"account_uuid" UUID REFERENCES accounts(uuid) NOT NULL,
	"user_uuid" UUID REFERENCES users(uuid) NOT NULL,
	"created_at" timestamp NOT NULL,
	"modified_at" timestamp NOT NULL,
	"deleted_at" timestamp,

	"processed" boolean NOT NULL,
	"iso_currency_code" varchar,
	"amount" varchar,
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
		return errors.Wrapf(err, "failed to ensure source_transactions table")
	}

	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON source_transactions USING btree(account_uuid)`)
	if err != nil {
		return errors.Wrap(err, "failed to ensure account_uuid index for source_transactions")
	}

	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON source_transactions USING btree(plaid_transaction_id)`)
	if err != nil {
		return errors.Wrap(err, "failed to ensure plaid_transaction_id index for source_transactions")
	}

	return nil
}

func (a *DBAgent) UpsertSourceTransaction(ctx context.Context, transaction SourceTransaction) (string, bool, error) {
	row := a.db.QueryRowContext(ctx, `
INSERT INTO "source_transactions" (
	"account_uuid",
	"user_uuid",
	"created_at",
	"modified_at",

	"processed",
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
	$1, $2, NOW(), NOW(),
	FALSE, $3, $4, $5,
	$6, $7, $8, $9, $10, $11, $12, $13
) ON CONFLICT ("uuid")
DO UPDATE SET
	"modified_at" = NOW(),
	"amount" = $5,
	"plaid_pending" = $9,
	"plaid_pending_transaction_id" = $10
RETURNING "uuid", "created_at" = "modified_at"`,
		transaction.AccountUUID,
		transaction.UserUUID,

		transaction.ISOCurrencyCode,
		transaction.AmountFloat(),
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
	var uuid string
	err := row.Scan(&uuid, &isNew)
	return uuid, isNew, errors.Wrapf(err, "failed to upsert to source_transactions table for plaid transaction %s", transaction.PlaidID)
}

func (a *DBAgent) DeleteSourceTransactionByPlaidID(ctx context.Context, plaidTransactionID string) error {
	_, err := a.db.ExecContext(ctx, `
UPDATE "source_transactions"
SET "deleted_at" = NOW()
WHERE "plaid_transaction_id" = $1`,
		plaidTransactionID,
	)
	return errors.Wrapf(err, "failed to insert into source_transactions table")
}

type SourceTransactionQuery struct {
	AccountUUID      string `json:"account_uuid"`
	IncludeProcessed bool   `json:"include_processed"`
}

const SourceTransactionQueryTemplate = `SELECT %s
FROM "source_transactions"
WHERE
	"deleted_at" IS NULL
	AND "user_uuid" = $1
	AND "account_uuid" = $2
	%s
`

func (q SourceTransactionQuery) whereClauseAddendum() string {
	if !q.IncludeProcessed {
		return `AND "processed" IS FALSE`
	}
	return ""
}

func (q SourceTransactionQuery) Name() string {
	return "SourceTransaction"
}
func (q SourceTransactionQuery) CountQuery() string {
	return fmt.Sprintf(
		SourceTransactionQueryTemplate,
		`COUNT(*)`,
		q.whereClauseAddendum(),
	)
}
func (q SourceTransactionQuery) Query() string {
	return fmt.Sprintf(
		SourceTransactionQueryTemplate,
		StandardSourceTransactionFieldNameList,
		q.whereClauseAddendum()+`
ORDER BY created_at
OFFSET $3 LIMIT $4
`,
	)
}
func (q SourceTransactionQuery) CountArgs(userUUID string) []interface{} {
	return []interface{}{userUUID, q.AccountUUID}
}
func (q SourceTransactionQuery) Args(userUUID string, skip int64) []interface{} {
	return []interface{}{userUUID, q.AccountUUID, skip, SourceTransactionMaxPageSize}
}

//TODO Are these helpers necessary, or should the DBAgent only have
// generic querying functionality, and then I write a bunch of
// implementations of Query? that'd give me the power to genericize
// the API side of things as well, since all query endpoints would be
// calling the same DBAgent methods
//
// verdict: probably don't do that yet? feels a little premature

func (a *DBAgent) SourceTransactionsQueryPreFlight(ctx context.Context, auth Authorization, q SourceTransactionQuery) (int64, error) {
	row := a.db.QueryRowContext(ctx, q.CountQuery(), q.CountArgs(auth.UserUUID)...)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get source_transaction count from table: %w", err)
	}

	return count, nil
}

func (a *DBAgent) StartSourceTransactionsQuery(ctx context.Context, auth Authorization, q SourceTransactionQuery) (ts []SourceTransaction, err error) {
	return a.queryHelper(ctx, auth, q, 0)
}

func (a *DBAgent) ContinueSourceTransactionsQuery(ctx context.Context, auth Authorization, q SourceTransactionQuery, skip int64) (ts []SourceTransaction, err error) {
	return a.queryHelper(ctx, auth, q, skip)
}

const SourceTransactionMaxPageSize = 20
