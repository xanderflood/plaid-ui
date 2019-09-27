package db

import (
	"database/sql"
)

//DB is the minimal database interface to back the app
//go:generate counterfeiter . DB
type DB interface {
	EnsureAccountsTable() error
	CreateAccount(acct Account) (string, error)
	// GetAccount(uuid string) (Account, error)
	GetAccounts() ([]Account, error)

	//these mark a given account as WebhookConfigured or not
	ConfigureAccount(uuid string) error
	DeconfigureAccount(uuid string) error

	// AddTransaction(transaction Transaction) (string, error)
}

//DBAgent implements DB using a *sql.DB
type DBAgent struct {
	db     *sql.DB
	uuider UUIDer
}

//NewDBAgent create a new DBAgent
func NewDBAgent(db *sql.DB) *DBAgent {
	return &DBAgent{
		db:     db,
		uuider: UUIDGenerator{},
	}
}
