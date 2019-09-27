package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/plaid/plaid-go/plaid"

	"github.com/xanderflood/plaid-ui/pkg/db"
	pkgplaid "github.com/xanderflood/plaid-ui/pkg/plaid"

	//postgres driver for db/sql
	_ "github.com/lib/pq"
)

// Fill with your Plaid API keys - https://dashboard.plaid.com/account/keys
var (
	PLAID_CLIENT_ID            = os.Getenv("PLAID_CLIENT_ID")
	PLAID_SECRET               = os.Getenv("PLAID_SECRET")
	PLAID_PUBLIC_KEY           = os.Getenv("PLAID_PUBLIC_KEY")
	APP_PORT                   = os.Getenv("APP_PORT")
	POSTGRES_CONNECTION_STRING = os.Getenv("POSTGRES_CONNECTION_STRING")
)

var plaidClientOptions = plaid.ClientOptions{
	ClientID:  PLAID_CLIENT_ID,
	Secret:    PLAID_SECRET,
	PublicKey: PLAID_PUBLIC_KEY,
	// Use 'sandbox' to test with fake credentials in Plaid's Sandbox environment
	// Use `development` to test with real credentials while developing
	// Use `production` to go live with real users
	Environment: plaid.Sandbox,

	HTTPClient: &http.Client{},
}

var plaidClient *plaid.Client
var dbClient db.DB

func init() {
	var err error

	plaidClient, err = plaid.NewClient(plaidClientOptions)
	if err != nil {
		log.Fatalf("couldn't initialize Plaid client: %s", err.Error())
	}

	sqlDB, err := sql.Open("postgres", POSTGRES_CONNECTION_STRING)
	if err != nil {
		log.Fatalf("couldn't initialize database connection: %s", err.Error())
	}

	dbClient = db.NewDBAgent(sqlDB)
	if err = dbClient.EnsureAccountsTable(); err != nil {
		log.Fatalf("couldn't initialize accounts table: %s", err.Error())
	}
}

func registerPlaidItem(c *gin.Context) {
	publicToken := c.PostForm("public_token")
	exchangeTokenResponse, err := plaidClient.ExchangePublicToken(publicToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getItemResponse, err := plaidClient.GetItem(exchangeTokenResponse.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getInstitutionResponse, err := plaidClient.GetInstitutionByIDWithOptions(
		getItemResponse.Item.InstitutionID,
		plaid.GetInstitutionByIDOptions{IncludeOptionalMetadata: true},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getAccountsResponse, err := plaidClient.GetAccounts(exchangeTokenResponse.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, acct := range getAccountsResponse.Accounts {
		//TODO enable the webhook for each account

		spew.Dump(getItemResponse.Item.Webhook)

		_, err := dbClient.CreateAccount(
			db.Account{
				PlaidAccessToken:    exchangeTokenResponse.AccessToken,
				PlaidAccountID:      acct.AccountID,
				PlaidAccountName:    acct.Name,
				PlaidAccountType:    pkgplaid.AccountType(acct.Type),
				PlaidAccountSubtype: pkgplaid.AccountSubtype(acct.Subtype),

				PlaidItemID:          getItemResponse.Item.ItemID,
				PlaidInstitutionName: getInstitutionResponse.Institution.Name,
				PlaidInstitutionURL:  getInstitutionResponse.Institution.URL,
				PlaidInstitutionLogo: getInstitutionResponse.Institution.Logo,
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": exchangeTokenResponse.AccessToken,
		"item_id":      getItemResponse.Item.ItemID,
	})
}

func accounts(c *gin.Context) {
	accounts, err := dbClient.GetAccounts()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

// TODO this should accept an item_id or account_uuid in the request body
//
// func createPublicToken(c *gin.Context) {
// 	// Create a one-time use public_token for the Item.
// 	// This public_token can be used to initialize Link in update mode for a user
// 	publicToken, err := plaidClient.CreatePublicToken(accessToken)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"public_token": publicToken,
// 	})
// }

func main() {
	if APP_PORT == "" {
		APP_PORT = "8000"
	}

	r := gin.Default()
	r.LoadHTMLFiles("templates/index.tmpl")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"plaid_environment": "sandbox", // Switch this environment
			"plaid_public_key":  PLAID_PUBLIC_KEY,
			"plaid_webhook_url": "TODO",
		})
	})

	r.POST("/set_access_token", registerPlaidItem)
	r.GET("/accounts", accounts)
	// r.GET("/create_public_token", createPublicToken)

	r.Run(":" + APP_PORT)
}
