package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/xanderflood/plaid-ui/lib/page"
)

//ErrBadToken indicates that an invalid pagination token has been provided
var ErrBadToken = errors.New("bad pagination token")

//DB is the minimal database interface to back the app
//go:generate counterfeiter . DB
type DB interface {
	EnsureUsersTable(ctx context.Context) error
	EnsureAccountsTable(ctx context.Context) error
	EnsureSourceTransactionsTable(ctx context.Context) error
	// EnsureTransactionsTable(ctx context.Context) error
	// EnsureCategoryTable(ctx context.Context) error
	// EnsureCategoryTransactionsTable(ctx context.Context) error

	//TODO add `auth Authorization` to every function, and use it
	// to build extra WHERE clauses

	RegisterUser(ctx context.Context, uuid string, email string) error
	CheckUser(ctx context.Context, uuid string) (bool, error)

	CreateAccount(ctx context.Context, userUUID string, acct Account) (string, error)
	GetAccountsByPlaidItemID(ctx context.Context, itemID string) ([]Account, error)
	GetAccounts(ctx context.Context, userUUID string) ([]Account, error)
	ConfigureAccount(ctx context.Context, userUUID string, uuid string) error
	DeconfigureAccount(ctx context.Context, userUUID string, uuid string) error

	UpsertSourceTransaction(ctx context.Context, transaction SourceTransaction) (string, bool, error)
	DeleteSourceTransactionByPlaidID(ctx context.Context, plaidTransactionID string) error
	StartSourceTransactionsQuery(ctx context.Context, auth Authorization, q SourceTransactionQuery) ([]SourceTransaction, string, error)
	ContinueSourceTransactionsQuery(ctx context.Context, auth Authorization, token string) ([]SourceTransaction, string, error)

	//first
	// TODO GetTransactions(ctx context.Context, userUUID string)
	// TODO DeleteTransaction(ctx context.Context, userUUID string, uuid string)
	// TODO GetAmortizedTransactionsForPeriod(ctx context.Context, userUUID string, start time.Time, end time.Time)

	//then
	// TODO ProcessSourceTransaction(
	//     ctx context.Context,
	//     sourceTransactionUUID string,
	//     transactions []Transaction,
	// )

	//now think about expanding the SPA

	//then
	// TODO CreateCategory(ctx context.Context, userUUID string, name string, deductibility *float64)
	// TODO UpdateCategory(ctx context.Context, userUUID string, name *string, deductibility *float64)
	// TODO GetCategory(ctx context.Context, userUUID string, uuid string)
	// TODO DeleteCategory(ctx context.Context, userUUID string, uuid string)

	//finally
	// TODO CategorizeTransaction
	// TODO UncategorizeTransaction

	//and then some transaction update functions
}

//DBAgent implements DB using a *sql.DB
type DBAgent struct {
	db      *sql.DB
	uuider  UUIDer
	tokener page.Tokener
}

//NewDBAgent create a new DBAgent
func NewDBAgent(db *sql.DB) *DBAgent {
	return &DBAgent{
		db:      db,
		uuider:  UUIDGenerator{},
		tokener: page.Base64JSONTokener{},
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
	err = db.EnsureSourceTransactionsTable(ctx)
	if err != nil {
		return err
	}
	// TODO
	// err = db.EnsureTransactionsTable(ctx)
	// if err != nil {
	// 	return err
	// }
	// err = db.EnsureCategoryTable(ctx)
	// if err != nil {
	// 	return err
	// }
	// err = db.EnsureCategoryTransactionsTable(ctx)
	// if err != nil {
	// 	return err
	// }
	return nil
}
