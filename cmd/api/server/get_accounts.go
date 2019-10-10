package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
