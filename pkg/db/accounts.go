package db

import (
	"context"

	"github.com/pkg/errors"
)

var ErrNoSuchAccount = errors.New("no such account")

//EnsureAccountsTable EnsureAccountsTable
func (a *DBAgent) EnsureAccountsTable(ctx context.Context) error {
	_, err := a.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS "accounts"
(	"uuid" UUID,
	"user_uuid" UUID REFERENCES users(uuid),
	"created_at" timestamp NOT NULL,
	"modified_at" timestamp NOT NULL,
	"deleted_at" timestamp,

	"webhook_configured" boolean DEFAULT false,

	"access_token" varchar NOT NULL,
	"plaid_account_id" varchar NOT NULL,
	"plaid_account_name" varchar NOT NULL,
	"plaid_account_type" varchar NOT NULL,
	"plaid_account_subtype" varchar NOT NULL,

	"plaid_item_id" varchar NOT NULL,
	"plaid_institution_name" varchar NOT NULL,
	"plaid_institution_url" varchar NOT NULL,
	"plaid_institution_logo" varchar NOT NULL,
	PRIMARY KEY ("uuid")
)`)
	if err != nil {
		return errors.Wrapf(err, "failed to ensure accounts table")
	}

	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON accounts USING btree(plaid_item_id)`)
	if err != nil {
		return errors.Wrap(err, "failed to ensure plaid_item_id index for accounts")
	}

	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON accounts USING btree(user_uuid)`)
	return errors.Wrap(err, "failed to ensure user_uuid index for accounts")

}

//CreateAccount inserts an account into the table
func (a *DBAgent) CreateAccount(ctx context.Context, userUUID string, acct Account) (string, error) {
	uuid := a.uuider.UUID()
	_, err := a.db.ExecContext(ctx, `
INSERT INTO "accounts" (
	"uuid",
	"user_uuid"
	"created_at",
	"modified_at",

	"access_token",
	"plaid_account_id",
	"plaid_account_name",
	"plaid_account_type",
	"plaid_account_subtype",

	"plaid_item_id",
	"plaid_institution_name",
	"plaid_institution_url",
	"plaid_institution_logo"
) VALUES (
	$1, $2, NOW(), NOW(),
	$3, $4, $5, $6, $7,
	$8, $9, $10, $11
)`,
		uuid,
		userUUID,
		acct.PlaidAccessToken,
		acct.PlaidAccountID,
		acct.PlaidAccountName,
		acct.PlaidAccountType,
		acct.PlaidAccountSubtype,
		acct.PlaidItemID,
		acct.PlaidInstitutionName,
		acct.PlaidInstitutionURL,
		acct.PlaidInstitutionLogo,
	)
	if err != nil {
		return "", errors.Wrapf(err, "failed to insert into accounts table")
	}
	return uuid, nil
}

//GetAccounts gets all the accounts
func (a *DBAgent) GetAccounts(ctx context.Context, userUUID string) ([]Account, error) {
	rows, err := a.db.QueryContext(ctx, `
SELECT
	"uuid",
	"created_at",
	"modified_at",

	"access_token",
	"plaid_account_id",
	"plaid_account_name",
	"plaid_account_type",
	"plaid_account_subtype",

	"plaid_item_id",
	"plaid_institution_name",
	"plaid_institution_url",
	"plaid_institution_logo"
FROM "accounts"
WHERE
	"deleted_at" IS NULL
	AND
	"user_uuid" = $1
`,
		userUUID,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get accounts from table")
	}

	var accounts []Account
	for rows.Next() {
		var account Account
		err = rows.Scan(
			&account.UUID,
			&account.CreatedAt,
			&account.ModifiedAt,
			&account.PlaidAccessToken,
			&account.PlaidAccountID,
			&account.PlaidAccountName,
			&account.PlaidAccountType,
			&account.PlaidAccountSubtype,
			&account.PlaidItemID,
			&account.PlaidInstitutionName,
			&account.PlaidInstitutionURL,
			&account.PlaidInstitutionLogo,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to scan result of querying for all accounts")
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

//ConfigureAccount mark an account as webhook-configured
func (a *DBAgent) ConfigureAccount(ctx context.Context, userUUID string, uuid string) error {
	return a.setAccountConfigured(ctx, userUUID, uuid, true)
}

//DeconfigureAccount mark an account as not webhook-configured
func (a *DBAgent) DeconfigureAccount(ctx context.Context, userUUID string, uuid string) error {
	return a.setAccountConfigured(ctx, userUUID, uuid, false)
}

//TODO make sure this can only access accounts owned by the user
func (a *DBAgent) setAccountConfigured(ctx context.Context, userUUID string, uuid string, val bool) error {
	_, err := a.db.ExecContext(ctx, userUUID, `
UPDATE "accounts"
SET "webhook_configured" = $1
WHERE "arrivals"."identifier" = $2`,
		val,
		uuid,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to update webhook_configured field for account `%s`", uuid)
	}
	return nil
}
