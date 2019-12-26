package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xanderflood/plaid-ui/lib/page"
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

//TODO how much can I genericize the logic around this query/pagination stuff?

//QuerySourceTransactions gets all the accounts
func (a ServerAgent) QuerySourceTransactions(c *gin.Context) {
	auth, ok := a.authorize(c)
	if !ok {
		return //an error response has already been sent
	}

	var req QuerySourceTransactionsRequest
	err := c.ShouldBindJSON(&req)
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
		query db.SourceTransactionQuery
		count int64
		skip  int64
	)
	if req.Token != "" {
		var td page.SkipTakeTokenData
		err := a.tokener.ParseToken(req.Token, &td)
		if err != nil {
			a.logger.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token string provided"})
			return
		}
		skip = td.Skip

		err = td.ParseQuery(&query)
		if err != nil {
			a.logger.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid query descriptor"})
			return
		}

		count, err = a.dbClient.SourceTransactionsQueryPreFlight(c, dbAuth, query)
		if err != nil {
			a.logger.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "source transactino pre-flight query failed"})
			return
		}

		ts, err = a.dbClient.ContinueSourceTransactionsQuery(c, dbAuth, query, td.Skip)
		if err != nil {
			a.logger.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed continuing source transaction query"})
			return
		}
	} else {
		query = db.SourceTransactionQuery{
			AccountUUID:      req.AccountUUID,
			IncludeProcessed: req.IncludeProcessed,
		}

		count, err = a.dbClient.SourceTransactionsQueryPreFlight(c, dbAuth, query)
		if err != nil {
			a.logger.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "source transactino pre-flight query failed"})
			return
		}

		ts, err = a.dbClient.StartSourceTransactionsQuery(c, dbAuth, query)
		if err != nil {
			a.logger.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed starting source transaction query"})
			return
		}
	}

	//build a next token, unless we have all the results
	var token string
	if skip+int64(len(ts)) < count {
		td := page.SkipTakeTokenData{Skip: skip + int64(len(ts))}
		td.SetQuery(query)
		tokenBs, err := a.tokener.ToTokenString(td)
		if err != nil {
			a.logger.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed generating next token"})
			return
		}

		token = string(tokenBs)
	}

	c.JSON(http.StatusOK, gin.H{
		"source_transactions": ts,
		"results":             count,
		"results_in_page":     len(ts),
		"next_token":          token,
	})
}
