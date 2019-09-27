package server

import (
	"net/url"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/xanderflood/plaid-ui/cmd/api/server/auth"
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
	jwtSigningSecret string
	hardcodedJWT     string

	plaidClient plaidapi.Client
	dbClient    db.DB
}

func (a ServerAgent) AddRoutes(
	e *gin.Engine,
) {
	e.GET("/", a.ServeSPA)

	//JWT endpoints
	jwtGroup := e.Group("/api/v1",
		auth.JWTMiddleware(
			a.jwtSigningSecret,
			&jwt.Parser{ValidMethods: []string{"HS256"}},
		),
	)
	jwtGroup.POST("/add_plaid_item", a.AddPlaidItem)
	jwtGroup.GET("/get_accounts", a.GetAccounts)
}

func NewServer(
	serviceDomain string,
	spaTermplateName string,
	plaidPublicKey string,
	plaidWebhookPath string,
	plaidEnvironment string,
	jwtSigningSecret string,
	hardcodedJWT string,

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
		jwtSigningSecret: jwtSigningSecret,
		hardcodedJWT:     hardcodedJWT,

		plaidClient: plaidClient,
		dbClient:    dbClient,
	}
}
