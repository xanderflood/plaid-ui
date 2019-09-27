package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//ServeSPA
func (a ServerAgent) ServeSPA(c *gin.Context) {
	c.HTML(http.StatusOK, a.spaTermplateName, gin.H{
		"plaid_environment": a.plaidEnvironment,
		"plaid_public_key":  a.plaidPublicKey,
		"plaid_webhook_url": a.plaidWebhookURL,
		"hardcoded_jwt":     a.hardcodedJWT,
	})
}
