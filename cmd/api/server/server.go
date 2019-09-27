package server

import (
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/xanderflood/plaid-ui/pkg/db"
	"github.com/xanderflood/plaid-ui/pkg/plaidapi"
)

//Server is the gin server interface for the public API
//go:generate counterfeiter . Server
type Server interface {
	ServeSPA(c *gin.Context)
	RegisterPlaidItem(c *gin.Context)
	GetAccounts(c *gin.Context)
	// GenericWebhook(c *gin.Context)
}

//ServerAgent implements Server
type ServerAgent struct {
	serviceDomain    string
	spaTermplateName string
	plaidPublicKey   string
	plaidWebhookURL  string
	plaidEnvironment string

	plaidClient plaidapi.Client
	dbClient    db.DB
}

//TODO ServerAgent#BuildGinServer
// Should build the appropriate hierarchy with a _certain_
// subset of endpoints wrapped in the middleware from the
// `auth` package.
// Also, wrap the webhook in ipfilter.

func NewServer(
	serviceDomain string,
	spaTermplateName string,
	plaidPublicKey string,
	plaidWebhookPath string,
	plaidEnvironment string,

	plaidClient plaidapi.Client,
	dbClient db.DB,
) ServerAgent {
	plaidWebhookURL := (&url.URL{
		Scheme: "https",
		Host:   serviceDomain,
		Path:   plaidWebhookPath,
	}).String()

	return ServerAgent{
		serviceDomain:    serviceDomain,
		spaTermplateName: spaTermplateName,
		plaidPublicKey:   plaidPublicKey,
		plaidWebhookURL:  plaidWebhookURL,
		plaidEnvironment: plaidEnvironment,

		plaidClient: plaidClient,
		dbClient:    dbClient,
	}
}
