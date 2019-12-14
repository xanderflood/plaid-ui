package db

// func (a *DBAgent) EnsureTransactionsTable(ctx context.Context) error {
// 	_, err := a.db.ExecContext(ctx, `
// CREATE TABLE IF NOT EXISTS "transactions"
// (	"uuid" UUID DEFAULT gen_random_uuid(),
// 	"account_uuid" UUID REFERENCES accounts(uuid),
// 	"user_uuid" UUID REFERENCES users(uuid),
// 	"created_at" timestamp NOT NULL,
// 	"modified_at" timestamp NOT NULL,
// 	"deleted_at" timestamp,

// 	"source_transaction_uuid" UUID REFERENCES source_transactions(uuid),
// 	"inverse_transaction_uuid" UUID REFERENCES transactions(uuid),

// 	"iso_currency_code" varchar,
// 	"amount" varchar,
// 	"date" date,
// 	"amortization_period" tsrange,

// 	"account_id" varchar,
// 	"name" varchar,
// 	"category_id" varchar,
// 	"account_owner" varchar,
// 	"transaction_type" varchar,
// 	PRIMARY KEY ("uuid")
// )`)
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to ensure transactions table")
// 	}

// 	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON transactions USING btree(account_uuid)`)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to ensure account_uuid index for transactions")
// 	}
// 	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON transactions USING btree(user_uuid)`)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to ensure user_uuid index for transactions")
// 	}
// 	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON transactions USING btree(source_transaction_uuid)`)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to ensure source_transaction_uuid index for source_transactions")
// 	}
// 	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON transactions USING gist(amortization_period)`)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to ensure amortization_period index for source_transactions")
// 	}

// 	return nil
// }

// TODO
// TODO
// TODO
// TODO before this, focus on pagination, and on creating a view
// TODO for source transactions and one for unprocessed source transactions
// TODO
// TODO
// TODO
// func (a *DBAgent) GetTransactions(ctx context.Context, userUUID string) ([]Transaction, error) {
// 	rows, err := a.db.QueryContext(ctx, `
// SELECT
// 	"uuid",
// 	"account_uuid",
// 	"user_uuid",
// 	"created_at",
// 	"modified_at",

// 	"iso_currency_code",
// 	"amount",
// 	"date",

// 	"plaid_account_id",
// 	"plaid_name",
// 	"plaid_category_id",
// 	"plaid_pending",
// 	"plaid_pending_transaction_id",
// 	"plaid_account_owner",
// 	"plaid_transaction_id",
// 	"plaid_type"
// FROM "accounts"
// WHERE
// 	"deleted_at" IS NULL
// 	AND "user_uuid" = $1
// 	AND "trans" = $2
// `,
// 		userUUID,
// 		accountUUID,
// 	)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to get source_transactions from table")
// 	}

// 	//
// 	//
// 	//
// 	//
// 	//
// }

// func (a *DBAgent) DeleteTransaction(ctx context.Context, userUUID string, uuid string) {

// }

// func (a *DBAgent) GetAmortizedTransactionsForPeriod(ctx context.Context, userUUID string, start time.Time, end time.Time) {

// }

// //TODO
// // TODO ProcessSourceTransaction(
// //     ctx context.Context,
// //     sourceTransactionUUID string,
// //     transactions []Transaction,
// // )
