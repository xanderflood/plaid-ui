package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	flag "github.com/jessevdk/go-flags"
	"github.com/plaid/plaid-go/plaid"

	"github.com/xanderflood/plaid-ui/cmd/api/server"
	"github.com/xanderflood/plaid-ui/cmd/api/server/auth"
	"github.com/xanderflood/plaid-ui/cmd/api/server/views"
	"github.com/xanderflood/plaid-ui/lib/tools"
	"github.com/xanderflood/plaid-ui/pkg/db"

	//postgres driver for db/sql
	_ "github.com/lib/pq"
)

var options struct {
	ServiceDomain            string `long:"service-domain"             env:"SERVICE_DOMAIN"              required:"true"`
	PlaidClientID            string `long:"plaid-client-id"            env:"PLAID_CLIENT_ID"             required:"true"`
	PlaidSecret              string `long:"plaid-secret"               env:"PLAID_SECRET"                required:"true"`
	PlaidPublicKey           string `long:"plaid-public-key"           env:"PLAID_PUBLIC_KEY"            required:"true"`
	PlaidEnvironment         string `long:"plaid-environment"          env:"PLAID_ENVIRONMENT"           required:"true"`
	PostgresConnectionString string `long:"postgres-connection-string" env:"POSTGRES_CONNECTION_STRING"  required:"true"`
	JWTSigningSecret         string `long:"jwt-signing-secret"         env:"JWT_SIGNING_SECRET"          required:"true"`
	LoginBaseURL             string `long:"login-base-url"             env:"LOGIN_BASE_URL"              required:"true"`

	Port  string `long:"port"          env:"PORT" default:"8000"`
	Debug bool   `long:"debug"         env:"DEBUG"`
}

func main() {
	_, err := flag.Parse(&options)
	if err != nil {
		log.Fatal(err)
	}

	plaidClient, err := plaid.NewClient(plaid.ClientOptions{
		ClientID:  options.PlaidClientID,
		Secret:    options.PlaidSecret,
		PublicKey: options.PlaidPublicKey,

		// Use 'sandbox' to test with fake credentials in Plaid's Sandbox environment
		// Use `development` to test with real credentials while developing
		// Use `production` to go live with real users
		Environment: plaid.Sandbox,

		HTTPClient: &http.Client{},
	})
	if err != nil {
		log.Fatalf("couldn't initialize Plaid client: %s", err.Error())
	}

	sqlDB, err := sql.Open("postgres", options.PostgresConnectionString)
	if err != nil {
		log.Fatalf("couldn't initialize database connection: %s", err.Error())
	}

	dbClient := db.NewDBAgent(sqlDB)
	if err = db.EnsureTables(context.Background(), dbClient); err != nil {
		log.Fatalf("couldn't initialize accounts table: %s", err.Error())
	}

	loginBaseURL, err := url.Parse(options.LoginBaseURL)
	if err != nil {
		log.Fatalf("login base URL `%s` was malformed: %s", options.LoginBaseURL, err.Error())
	}

	logger := tools.NewStdoutLogger()

	renderer := views.NewRenderer(
		logger,
		options.PlaidEnvironment,
		options.PlaidPublicKey,
		"/v1/plaid/webhook", //TODO refactor the router logic
		map[views.TemplateName]string{
			views.TemplateNameSPA:           "index.tmpl",
			views.TemplateNameNotRegistered: "not_registered.tmpl",
			views.TemplateNameErrorCode:     "error_code.tmpl",
		},
	)

	authMgr := auth.NewAuthorizationManager(
		logger,
		renderer,
		options.JWTSigningSecret,
		&jwt.Parser{ValidMethods: []string{"HS256"}},
		dbClient,
		loginBaseURL,
	)

	srv := server.NewServer(
		logger,

		options.ServiceDomain,
		"/v1/plaid/webhook", //TODO refactor the router logic

		authMgr,
		renderer,
		auth.GetAuthorizationFromContext,
		plaidClient,
		dbClient,
	)

	//build the gin server
	r := gin.Default()

	r.LoadHTMLFiles(
		"templates/index.tmpl",
		"templates/not_registered.tmpl",
		"templates/error_code.tmpl",
	)
	r.Static("/static", "./static")
	server.AddRoutes(r, srv)

	r.Run(":" + options.Port)
}
