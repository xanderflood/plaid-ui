package views

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xanderflood/plaid-ui/lib/tools"
)

type Renderer interface {
	RenderSPA(c *gin.Context)
	RenderNotRegistered(email string, c *gin.Context)
	RenderStatusCode(status int, message string, c *gin.Context)
}

type RendererAgent struct {
	logger           tools.Logger
	templateNames    map[TemplateName]string
	plaidEnvironment string
	plaidPublicKey   string
	plaidWebhookURL  string
}

type TemplateName string

const (
	TemplateNameSPA           TemplateName = "SPA"
	TemplateNameNotRegistered TemplateName = "NotRegistered"
	TemplateNameErrorCode     TemplateName = "404"
)

func NewRenderer(
	logger tools.Logger,
	plaidEnvironment string,
	plaidPublicKey string,
	plaidWebhookURL string,
	templateNames map[TemplateName]string,
) RendererAgent {
	return RendererAgent{
		logger:           logger,
		plaidEnvironment: plaidEnvironment,
		plaidPublicKey:   plaidPublicKey,
		plaidWebhookURL:  plaidWebhookURL,
		templateNames:    templateNames,
	}
}

//ServeSPA renders the single-page app
func (a RendererAgent) RenderSPA(c *gin.Context) {
	c.HTML(http.StatusOK, a.templateNames[TemplateNameSPA], gin.H{
		"plaid_environment": a.plaidEnvironment,
		"plaid_public_key":  a.plaidPublicKey,
		"plaid_webhook_url": a.plaidWebhookURL,
	})
	c.Abort()
}

//RenderNotRegistered serves the user-not-registered view
func (a RendererAgent) RenderNotRegistered(email string, c *gin.Context) {
	a.logger.Infof("%v %v", http.StatusForbidden, email)
	c.HTML(http.StatusForbidden, a.templateNames[TemplateNameNotRegistered], gin.H{
		"email": email,
	})
	c.Abort()
}

//RenderStatusCode serves the generic error page
func (a RendererAgent) RenderStatusCode(status int, message string, c *gin.Context) {
	a.logger.Infof("%v %v", status, message)
	c.HTML(status, a.templateNames[TemplateNameErrorCode], gin.H{
		"status_code": status,
		"message":     message,
	})
	c.Abort()
}
