package db

import (
	"time"

	"github.com/xanderflood/plaid-ui/pkg/plaidapi"
)

type Model struct {
	UUID       string     `json:"uuid"`
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt time.Time  `json:"modified_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

//Account represents a single bank account
type Account struct {
	Model

	WebhookConfigured bool `json:"webhook_configured"`

	PlaidAccessToken    string                  `json:"plaid_access_token"`
	PlaidAccountID      string                  `json:"plaid_account_id"`
	PlaidAccountName    string                  `json:"plaid_account_name"`
	PlaidAccountType    plaidapi.AccountType    `json:"plaid_account_type"`
	PlaidAccountSubtype plaidapi.AccountSubtype `json:"plaid_account_subtype"`

	PlaidItemID          string `json:"plaid_item_id"`
	PlaidInstitutionName string `json:"plaid_institution_name"`
	PlaidInstitutionURL  string `json:"plaid_institution_url"`
	PlaidInstitutionLogo string `json:"plaid_institution_logo"`
}
