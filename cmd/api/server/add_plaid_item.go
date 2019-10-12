package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plaid/plaid-go/plaid"
	"github.com/xanderflood/plaid-ui/pkg/db"
	"github.com/xanderflood/plaid-ui/pkg/plaidapi"
)

//AddPlaidItem adds all the accounts associated with this plaid item
func (a ServerAgent) AddPlaidItem(c *gin.Context) {
	authorization, ok := a.authorize(c)
	if !ok {
		return
	}

	publicToken := c.PostForm("public_token")
	exchangeTokenResponse, err := a.plaidClient.ExchangePublicToken(publicToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getItemResponse, err := a.plaidClient.GetItem(exchangeTokenResponse.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getInstitutionResponse, err := a.plaidClient.GetInstitutionByIDWithOptions(
		getItemResponse.Item.InstitutionID,
		plaid.GetInstitutionByIDOptions{IncludeOptionalMetadata: true},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getAccountsResponse, err := a.plaidClient.GetAccounts(exchangeTokenResponse.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, acct := range getAccountsResponse.Accounts {
		//TODO enable the webhook for each account
		//  is this going to be done automatically by the frontend?

		_, err := a.dbClient.CreateAccount(c,
			authorization.UserUUID,
			db.Account{
				PlaidAccessToken:    exchangeTokenResponse.AccessToken,
				PlaidAccountID:      acct.AccountID,
				PlaidAccountName:    acct.Name,
				PlaidAccountType:    plaidapi.AccountType(acct.Type),
				PlaidAccountSubtype: plaidapi.AccountSubtype(acct.Subtype),

				PlaidItemID:          getItemResponse.Item.ItemID,
				PlaidInstitutionName: getInstitutionResponse.Institution.Name,
				PlaidInstitutionURL:  getInstitutionResponse.Institution.URL,
				PlaidInstitutionLogo: getInstitutionResponse.Institution.Logo,
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": exchangeTokenResponse.AccessToken,
		"item_id":      getItemResponse.Item.ItemID,
	})
}
