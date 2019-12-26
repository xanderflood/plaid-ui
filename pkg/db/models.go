package db

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/xanderflood/plaid-ui/pkg/plaidapi"
)

//Model contains generic fields shared by all models
type Model struct {
	UUID       string     `json:"uuid"`
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt time.Time  `json:"modified_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

//Account represents a single bank account
type Account struct {
	Model

	UserUUID string `json:"user_uuid"`

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

const StandardAccountFieldNameList = `
	"uuid",
	"user_uuid",
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
`

func (a *Account) StandardFieldPointers() []interface{} {
	return []interface{}{
		&a.UUID,
		&a.UserUUID,
		&a.CreatedAt,
		&a.ModifiedAt,

		&a.PlaidAccessToken,
		&a.PlaidAccountID,
		&a.PlaidAccountName,
		&a.PlaidAccountType,
		&a.PlaidAccountSubtype,
		&a.PlaidItemID,
		&a.PlaidInstitutionName,
		&a.PlaidInstitutionURL,
		&a.PlaidInstitutionLogo,
	}
}

//SourceTransaction represents a single transaction from a source like Plaid
type SourceTransaction struct {
	Model

	AccountUUID string `json:"account_uuid"`
	UserUUID    string `json:"user_uuid"`

	Processed bool `json:"processed"`

	//TODO will NULL amounts be handled properly here?
	ISOCurrencyCode string       `json:"iso_currency_code"`
	Amount          *json.Number `json:"amount"`
	Date            string       `json:"date"`

	PlaidAccountID            string `json:"plaid_account_id"`
	PlaidName                 string `json:"plaid_name"`
	PlaidCategoryID           string `json:"plaid_category_id"`
	PlaidPending              bool   `json:"plaid_pending"`
	PlaidPendingTransactionID string `json:"plaid_pending_transaction_id"`
	PlaidAccountOwner         string `json:"plaid_account_owner"`
	PlaidID                   string `json:"plaid_transaction_id"`
	PlaidType                 string `json:"plaid_transaction_type"`
}

func (t SourceTransaction) AmountFloat() float64 {
	fl, _ := t.Amount.Float64()
	return fl
}

const StandardSourceTransactionFieldNameList = `
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
`

func (t *SourceTransaction) StandardFieldPointers() []interface{} {
	return []interface{}{
		&t.UUID,
		&t.AccountUUID,
		&t.UserUUID,
		&t.CreatedAt,
		&t.ModifiedAt,

		&t.ISOCurrencyCode,
		&t.Amount,
		&t.Date,

		&t.PlaidAccountID,
		&t.PlaidName,
		&t.PlaidCategoryID,
		&t.PlaidPending,
		&t.PlaidPendingTransactionID,
		&t.PlaidAccountOwner,
		&t.PlaidID,
		&t.PlaidType,
	}
}

//Transaction represents a single user-created transaction
type Transaction struct {
	Model

	AccountUUID           string `json:"account_uuid"`
	UserUUID              string `json:"user_uuid"`
	SourceTransactionUUID string `json:"source_transaction_uuid"`
	InverseOf             string `json:"inverse_transaction_uuid"`

	ISOCurrencyCode    string     `json:"iso_currency_code"`
	Amount             *big.Float `json:"amount"`
	Date               string     `json:"date"`
	AmortizationPeriod string     `json:"amortization_period"`

	AccountID    string `json:"account_id"`
	Name         string `json:"name"`
	CategoryID   string `json:"category_id"`
	AccountOwner string `json:"account_owner"`
	Type         string `json:"transaction_type"`
}

type Category struct {
	Model

	UserUUID string `json:"user_uuid"`
	Name     string `json:"name"`

	Deductibility *json.Number `json:"deductibility,omitempty"`

	//TODO add budgets to a category
	// or maybe BUDGET is a separate model that provides a rule
	// for mapping category/period pairs to amounts
}
