package server

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/xanderflood/plaid-ui/pkg/db"
)

//QuerySourceTransactionsRequest encodes a single request to query transactions
type QuerySourceTransactionsRequest struct {
	AccountUUID      string `json:"account_uuid"`
	IncludeProcessed bool   `json:"include_processed"`

	Token string `json:"token"`
}

//QuerySourceTransactionsResponse encodes a single response to query transactions
type QuerySourceTransactionsResponse struct {
	SourceTransactions []db.SourceTransaction `json:"source_transactions"`
	Token              string                 `json:"token,omitempty"`
}

//QuerySourceTransactions gets all the accounts
func (a ServerAgent) QuerySourceTransactions(c *gin.Context) {
	auth, ok := a.authorize(c)
	if !ok {
		return //an error response has already been generated
	}

	//
	//
	// TODO TODO TODO TODO
	// sign/encrypt all tokens so that the user_uuid can't be
	// tampered with
	//
	//
	// OR SOMETHING?
	//  e.g. add userUUID back to theh query functions and have the
	//   db client validate permissions?
	//
	//

	var req QuerySourceTransactionsRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		//TODO is this error message safe to expose?
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	spew.Dump(req)
	var (
		ts    []db.SourceTransaction
		token string
	)
	if req.Token != "" {
		ts, token, err = a.dbClient.ContinueSourceTransactionsQuery(c, req.Token)
	} else {
		query := db.SourceTransactionQuery{
			UserUUID:         auth.UserUUID,
			AccountUUID:      req.AccountUUID,
			IncludeProcessed: req.IncludeProcessed,
		}

		ts, token, err = a.dbClient.StartSourceTransactionsQuery(c, query)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"source_transactions": ts,
		"next_token":          token,
	})
}
