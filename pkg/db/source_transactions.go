package db

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/xanderflood/plaid-ui/lib/page"
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

func (q SourceTransactionQuery) Query() string {
	var addendum string
	if !q.IncludeProcessed {
		addendum = `AND "processed" IS FALSE`
	}
	return fmt.Sprintf(`SELECT %s
FROM "source_transactions"
WHERE
	"deleted_at" IS NULL
	AND "user_uuid" = $1
	AND "account_uuid" = $2
	%s
OFFSET $3 LIMIT $4
`, StandardSourceTransactionFieldNameList, addendum)
}
func (q SourceTransactionQuery) Args(userUUID string, skip int64) []interface{} {
	return []interface{}{userUUID, q.AccountUUID, skip, SourceTransactionMaxPageSize}
}

func (a *DBAgent) StartSourceTransactionsQuery(ctx context.Context, auth Authorization, q SourceTransactionQuery) ([]SourceTransaction, string, error) {
	ts, token, err := a.sourceTransactionQueryHelper(ctx, auth, 0, q)
	if err != nil {
		return nil, "", fmt.Errorf("failed to start query on source transactions: %w", err)
	}

	return ts, token, err
}

func (a *DBAgent) ContinueSourceTransactionsQuery(ctx context.Context, auth Authorization, token string) ([]SourceTransaction, string, error) {
	var td page.SkipTakeTokenData
	err := a.tokener.ParseToken(token, td)
	if err != nil {
		return nil, "", fmt.Errorf("invalid token string provided: %w", err)
	}

	var q SourceTransactionQuery
	err = td.ParseQuery(q)
	if err != nil {
		return nil, "", fmt.Errorf("invalid query descriptor provided: %w", err)
	}

	ts, token, err := a.sourceTransactionQueryHelper(ctx, auth, td.Skip, q)
	if err != nil {
		return nil, "", fmt.Errorf("failed to confinue query on source transactions with token %s: %w", token, err)
	}

	return ts, token, nil
}

const SourceTransactionMaxPageSize = 20

func (a *DBAgent) sourceTransactionQueryHelper(ctx context.Context, auth Authorization, skip int64, q SourceTransactionQuery) ([]SourceTransaction, string, error) {
	spew.Dump(q.Query(), q.Args(auth.UserUUID, skip))
	rows, err := a.db.QueryContext(ctx, q.Query(), q.Args(auth.UserUUID, skip)...)
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to get source_transactions from table")
	}

	var sourceTransactions []SourceTransaction
	for rows.Next() {
		var sourceTransaction SourceTransaction
		err = rows.Scan((&sourceTransaction).StandardFieldPointers()...)
		if err != nil {
			//TODO eliminate errors.Wrapf
			//TODOFIRST find out if fmt.Errorf handles nils the same as errors.Wrapf
			return nil, "", errors.Wrapf(err, "failed to scan result of querying for source_transactions")
		}

		sourceTransactions = append(sourceTransactions, sourceTransaction)
	}

	//if it's an empty page, we hit the end
	if len(sourceTransactions) == 0 {
		return nil, "", nil
	}

	td := page.SkipTakeTokenData{Skip: skip + int64(len(sourceTransactions))}
	td.SetQuery(q)
	tokenBs, err := a.tokener.ToTokenString(td)
	return sourceTransactions, string(tokenBs), err
}
