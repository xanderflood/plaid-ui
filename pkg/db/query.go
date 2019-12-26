package db

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

//Query supplies the SQL material necessary for a paginated query
type Query interface {
	Name() string
	CountQuery() string
	Query() string
	CountArgs(userUUID string) []interface{}
	Args(userUUID string, skip int64) []interface{}
}

func (a *DBAgent) queryHelper(ctx context.Context, auth Authorization, q Query, skip int64) ([]SourceTransaction, error) {
	spew.Dump(q.Query(), q.Args(auth.UserUUID, skip))
	rows, err := a.db.QueryContext(ctx, q.Query(), q.Args(auth.UserUUID, skip)...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s query: %w", q.Name(), err)
	}

	var sourceTransactions []SourceTransaction
	for rows.Next() {
		var sourceTransaction SourceTransaction
		err = rows.Scan((&sourceTransaction).StandardFieldPointers()...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan result of %s query: %w", q.Name(), err)
		}

		sourceTransactions = append(sourceTransactions, sourceTransaction)
	}

	//if it's an empty page, we hit the end
	if len(sourceTransactions) == 0 {
		return nil, nil
	}

	return sourceTransactions, err
}
