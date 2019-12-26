package server

import (
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/xanderflood/plaid-ui/cmd/api/server/auth"
	"github.com/xanderflood/plaid-ui/cmd/api/server/views"
	"github.com/xanderflood/plaid-ui/lib/page"
	"github.com/xanderflood/plaid-ui/lib/tools"
	"github.com/xanderflood/plaid-ui/pkg/db"
	"github.com/xanderflood/plaid-ui/pkg/plaidapi"
)

//Server is the gin server interface for the public API
//go:generate counterfeiter . Server
type Server interface {
	// frontend
	ServeSPA(c *gin.Context)

	// user api
	AddPlaidItem(c *gin.Context)
	GetAccounts(c *gin.Context)
	QuerySourceTransactions(c *gin.Context)

	// admin api
	RegisterUser(c *gin.Context)

	// plaid webhooks
	GenericPlaidWebhook(c *gin.Context)

	// authorization code
	BackendAuthorizationMiddleware(c *gin.Context)
	FrontendAuthorizationMiddleware(c *gin.Context)
}

//ServerAgent implements Server
type ServerAgent struct {
	logger tools.Logger

	serviceDomain   string
	plaidWebhookURL string

	authorize   auth.Getter
	renderer    views.Renderer
	plaidClient plaidapi.Client
	dbClient    db.DB
	tokener     page.Tokener

	backendJWTMiddleware  gin.HandlerFunc
	frontendJWTMiddleware gin.HandlerFunc
}

//AddRoutes accepts a *gin.Engine and adds all the
//necessary routes to it for this API.
func AddRoutes(e *gin.Engine, a Server) {
	frontend := e.Group("/", a.FrontendAuthorizationMiddleware)
	frontend.GET("/", a.ServeSPA)

	//webhook
	webhook := e.Group("/webhook")
	webhook.POST("/v1/plaid", a.GenericPlaidWebhook)

	//JWT endpoints
	backend := e.Group("/api/v1", a.BackendAuthorizationMiddleware)
	backend.POST("/add_plaid_item", a.AddPlaidItem)
	backend.GET("/get_accounts", a.GetAccounts)
	backend.POST("/query_source_transactions", a.QuerySourceTransactions)

	//admin endpoints
	adminGroup := backend.Group("/admin")
	adminGroup.POST("/register_user", a.RegisterUser)
}

//NewServer creates a new Server.
func NewServer(
	logger tools.Logger,

	serviceDomain string,

	authMgr auth.AuthorizationManager,
	renderer views.Renderer,
	authorize auth.Getter,
	plaidClient plaidapi.Client,
	dbClient db.DB,
) ServerAgent {
	plaidWebhookURL := (&url.URL{
		Scheme: "https",
		Host:   serviceDomain,
	}).String()

	return ServerAgent{
		logger: logger,

		serviceDomain:   serviceDomain,
		plaidWebhookURL: plaidWebhookURL,

		authorize:   authorize,
		renderer:    renderer,
		plaidClient: plaidClient,
		dbClient:    dbClient,
		tokener:     page.Base64JSONTokener{},

		backendJWTMiddleware:  authMgr.BackendMiddleware(),
		frontendJWTMiddleware: authMgr.FrontendMiddleware(),
	}
}

//BackendAuthorizationMiddleware callback for the backend authorization middleware
func (a ServerAgent) BackendAuthorizationMiddleware(c *gin.Context) {
	a.backendJWTMiddleware(c)
}

//FrontendAuthorizationMiddleware callback for the frontend authorization middleware
func (a ServerAgent) FrontendAuthorizationMiddleware(c *gin.Context) {
	a.frontendJWTMiddleware(c)
}
