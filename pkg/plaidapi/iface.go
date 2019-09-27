package plaidapi

import "github.com/plaid/plaid-go/plaid"

//Client isolates the necessary interface with a *plaid.Client
//go:generate counterfeiter . Client
type Client interface {
	ExchangePublicToken(publicToken string) (resp plaid.ExchangePublicTokenResponse, err error)
	GetItem(accessToken string) (resp plaid.GetItemResponse, err error)
	GetInstitutionByIDWithOptions(id string, options plaid.GetInstitutionByIDOptions) (resp plaid.GetInstitutionByIDResponse, err error)
	GetAccounts(accessToken string) (resp plaid.GetAccountsResponse, err error)
}
