package server

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/xanderflood/plaid-ui/pkg/db"
)

//QuerySourceTransactionsRequest encodes a single request to query transactions
type QuerySourceTransactionsRequest struct {
	UserUUID         string `json:"user_uuid" binding:"required"`
	AccountUUID      string `json:"account_uuid" binding:"required"`
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
		return //an error response has already been sent
	}

	var req QuerySourceTransactionsRequest
	err := c.ShouldBindJSON(&req)
	spew.Dump(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbAuth, err := auth.GetDBAuthorization(req.UserUUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var (
		ts    []db.SourceTransaction
		token string
	)
	if req.Token != "" {
		ts, token, err = a.dbClient.ContinueSourceTransactionsQuery(c, dbAuth, req.Token)
	} else {
		query := db.SourceTransactionQuery{
			AccountUUID:      req.AccountUUID,
			IncludeProcessed: req.IncludeProcessed,
		}

		ts, token, err = a.dbClient.StartSourceTransactionsQuery(c, dbAuth, query)
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
