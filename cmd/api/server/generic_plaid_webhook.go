package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xanderflood/plaid-ui/pkg/db"
	"github.com/xanderflood/plaid-ui/pkg/plaidapi"
)

type WebhookType string

const (
	ItemWebhookType         WebhookType = "ITEM"
	TransactionsWebhookType WebhookType = "TRANSACTIONS"
)

type WebhookCode string

const (
	InitialUpdate       WebhookCode = "INITIAL_UPDATE"
	HistoricalUpdate    WebhookCode = "HISTORICAL_UPDATE"
	DefaultUpdate       WebhookCode = "DEFAULT_UPDATE"
	TransactionsRemoved WebhookCode = "TRANSACTIONS_REMOVED"

	ItemWebhookUpdateAcknowledged WebhookCode = "WEBHOOK_UPDATE_ACKNOWLEDGED"
	ItemError                     WebhookCode = "ERROR"
)

func (t WebhookCode) IsRemoval() bool {
	return t == TransactionsRemoved
}

type WebhookRequest struct {
	Type                WebhookType     `json:"webhook_type"`
	Code                WebhookCode     `json:"webhook_code"`
	ItemID              string          `json:"item_id"`
	Error               json.RawMessage `json:"error"`
	Newtransactions     int             `json:"new_transactions"`
	RemovedTransactions []string        `json:"removed_transactions"`
	NewWebhookURL       string          `json:"new_webhook_url"`
}

//GenericPlaidWebhook accepts all Plaid webhook requests
func (a ServerAgent) GenericPlaidWebhook(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}

	var wr WebhookRequest
	err = json.Unmarshal(reqBody, &wr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	a.logger.Debugf("processing webhook - type: %s, code: %s, item ID: %s", wr.Type, wr.Code, wr.ItemID)

	if len(wr.ItemID) == 0 {
		return
	}

	accounts, err := a.dbClient.GetAccountsByPlaidItemID(c, wr.ItemID)
	if err != nil {
		a.logger.Errorf("failed getting accounts for plaid item `%s`: %s", wr.ItemID, err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(accounts) == 0 {
		a.logger.Debugf("received webhook request for unrecognized item_id", wr.ItemID)
	}

	switch wr.Type {
	case ItemWebhookType:
		switch wr.Code {
		case ItemWebhookUpdateAcknowledged:
			if wr.NewWebhookURL != a.plaidWebhookURL {
				_, err := a.plaidClient.UpdateItemWebhook(accounts[0].PlaidItemID, a.plaidWebhookURL)
				if err != nil {
					a.logger.Errorf("failed processing webhook-update webhook for plaid item `%s` with `%s` as value: %s", wr.ItemID, wr.NewWebhookURL, err.Error())
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			}
			return

		case ItemError:
			a.logger.Errorf("received an error webhook from Plaid: %s: %s", wr.Error, err.Error())
			return

		default:
			a.logger.Errorf("invalid transaction webhook code `%s`: %s", wr.Code, err.Error())
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	case TransactionsWebhookType:
		a.logger.Infof("processing transaction webhook code `%s` for item `%s`", wr.Code, wr.ItemID)
		switch wr.Code {
		case InitialUpdate:
			_, err := a.transactionWebhookAddHelper(c, wr.Newtransactions, wr.ItemID)
			if err != nil {
				a.logger.Errorf("failed processing transaction webhook for plaid item `%s` with `%v` items: %s", wr.ItemID, wr.Newtransactions, err.Error())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			return

		case HistoricalUpdate:
			_, err := a.transactionWebhookAddHelper(c, wr.Newtransactions, wr.ItemID)
			if err != nil {
				a.logger.Errorf("failed processing transaction webhook for plaid item `%s` with `%v` items: %s", wr.ItemID, wr.Newtransactions, err.Error())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			return

		case DefaultUpdate:
			_, err := a.transactionWebhookAddHelper(c, wr.Newtransactions, wr.ItemID)
			if err != nil {
				a.logger.Errorf("failed processing transaction webhook for plaid item `%s` with `%v` items: %s", wr.ItemID, wr.Newtransactions, err.Error())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			return

		case TransactionsRemoved:
			for _, tid := range wr.RemovedTransactions {
				err := a.dbClient.DeleteTransactionByPlaidID(c, tid)
				if err != nil {
					a.logger.Errorf("failed processing transaction removal webhook: %s", err.Error())
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			}
			return

		default:
			a.logger.Errorf("invalid transaction webhook code `%s`: %s", wr.Code, err.Error())
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	default:
		//do nothing
		return
	}
}

func (a ServerAgent) addLimitedTransactionsForDate(ctx context.Context, date time.Time, remaining *int, userUUID string, accessToken string, accounts map[string]db.Account) error {
	//TODO I've been getting 429s from this - does it have built-in retry logic?
	getTransactionsResp, err := a.plaidClient.GetTransactions(accessToken, date.Format(plaidapi.DateFormat), date.Format(plaidapi.DateFormat))
	if err != nil {
		return err
	}

	for _, plaidTransaction := range getTransactionsResp.Transactions {
		transaction := db.Transaction{
			AccountUUID: accounts[plaidTransaction.AccountID].UUID,
			UserUUID:    userUUID,

			ISOCurrencyCode: plaidTransaction.ISOCurrencyCode,
			Amount:          big.NewFloat(plaidTransaction.Amount),
			Date:            plaidTransaction.Date,

			PlaidAccountID:            plaidTransaction.AccountID,
			PlaidName:                 plaidTransaction.Name,
			PlaidCategoryID:           plaidTransaction.CategoryID,
			PlaidPending:              plaidTransaction.Pending,
			PlaidPendingTransactionID: plaidTransaction.PendingTransactionID,
			PlaidAccountOwner:         plaidTransaction.AccountOwner,
			PlaidID:                   plaidTransaction.ID,
			PlaidType:                 plaidTransaction.Type,
		}

		_, isNew, err := a.dbClient.UpsertTransaction(ctx, transaction)
		if err != nil {
			return err
		}
		if isNew {
			*remaining--

			if *remaining <= 0 {
				return nil
			}
		}
	}

	return nil
}

func (a ServerAgent) transactionWebhookAddHelper(ctx context.Context, max int, itemID string) (int, error) {
	//start with tomorrow to avoid timezone issues
	date := time.Now().Truncate(time.Hour * 24).Add(time.Hour * 24)

	accts, err := a.dbClient.GetAccountsByPlaidItemID(ctx, itemID)
	if err != nil {
		return 0, err
	}

	if len(accts) == 0 {
		return 0, fmt.Errorf("itemID `%s` has no plaid accounts", itemID)
	}

	// TODO: if the webhook is hit before the add_plaid_item call completes,
	//  the account record won't exist. To address that, we could have the
	//  webhook simply queue a re-tryable job instead of doing the upserts
	//  right away.
	// However this might not be a real issue if Plaid makes a habit of delaying a few seconds :shrug:

	var accountMapping = map[string]db.Account{}
	var accessToken = accts[0].PlaidAccessToken
	var userUUID string
	for _, account := range accts {
		userUUID = account.UserUUID // these will all be the same value
		accountMapping[account.PlaidAccountID] = account
	}

	var remaining = max
	for remaining > 0 {
		// TODO Put some reasonable limits so that it doesn't go crazy if a transaction gets sent twice

		err := a.addLimitedTransactionsForDate(ctx, date, &remaining, userUUID, accessToken, accountMapping)
		if err != nil {
			return remaining, err
		}

		date = date.Add(-time.Hour * 24)
	}

	return remaining, nil
}
