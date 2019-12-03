package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

//ErrBadToken indicates that an invalid pagination token has been provided
var ErrBadToken = errors.New("bad pagination token")

//DB is the minimal database interface to back the app
//go:generate counterfeiter . DB
type DB interface {
	EnsureUsersTable(ctx context.Context) error
	EnsureAccountsTable(ctx context.Context) error
	EnsureTransactionsTable(ctx context.Context) error

	RegisterUser(ctx context.Context, uuid string, email string) error
	CheckUser(ctx context.Context, uuid string) (bool, error)

	CreateAccount(ctx context.Context, userUUID string, acct Account) (string, error)
	GetAccountsByPlaidItemID(ctx context.Context, itemID string) ([]Account, error)
	GetAccounts(ctx context.Context, userUUID string) ([]Account, error)
	ConfigureAccount(ctx context.Context, userUUID string, uuid string) error
	DeconfigureAccount(ctx context.Context, userUUID string, uuid string) error

	UpsertTransaction(ctx context.Context, transaction Transaction) (string, bool, error)
	DeleteTransactionByPlaidID(ctx context.Context, plaidTransactionID string) error
	GetTransactions(ctx context.Context, accountUUID string, token string) ([]Transaction, error)
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

//EnsureTables builds out all the tables in order
func EnsureTables(ctx context.Context, db DB) error {
	err := db.EnsureUsersTable(ctx)
	if err != nil {
		return err
	}
	err = db.EnsureAccountsTable(ctx)
	if err != nil {
		return err
	}
	err = db.EnsureTransactionsTable(ctx)
	if err != nil {
		return err
	}
	return nil
}
