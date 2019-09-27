package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//GetAccounts gets all the accounts
func (a ServerAgent) GetAccounts(c *gin.Context) {
	accounts, err := a.dbClient.GetAccounts()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}
