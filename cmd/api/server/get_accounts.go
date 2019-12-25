package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

//TODO start using this request body so and the db.Authorization
//QueryAccountsRequest encodes a single request to query transactions
type QueryAccountsRequest struct {
	UserUUID string `json:"user_uuid"`

	//TODO add pagination
	Token string `json:"token"`
}

func (r QueryAccountsRequest) Validate() error {
	if r.UserUUID != "" {
		return errors.New("field `user_uuid` must be present and nonempty")
	}
	return nil
}

//GetAccounts gets all the accounts
func (a ServerAgent) GetAccounts(c *gin.Context) {
	auth, ok := a.authorize(c)
	if !ok {
		return //an error response has already been generated
	}

	accounts, err := a.dbClient.GetAccounts(c, auth.UserUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}
