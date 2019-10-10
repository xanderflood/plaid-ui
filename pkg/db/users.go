package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

//EnsureUsersTable EnsureUsersTable
func (a *DBAgent) EnsureUsersTable(ctx context.Context) error {
	_, err := a.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS "users"
(	"uuid" UUID,
	"email" varchar,
	"created_at" timestamp NOT NULL,
	"modified_at" timestamp NOT NULL,
	"deleted_at" timestamp,
	PRIMARY KEY ("uuid")
)`)
	if err != nil {
		return errors.Wrapf(err, "failed to ensure users table")
	}
	return nil
}

//RegisterUser whitelists a user for this service
func (a *DBAgent) RegisterUser(ctx context.Context, uuid string, email string) error {
	_, err := a.db.ExecContext(ctx, `
INSERT INTO "users" (
	"uuid",
	"email",
	"created_at",
	"modified_at"
) VALUES (
	$1, $2, NOW(), NOW()
)`,
		uuid, email,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to insert into users table")
	}
	return nil
}

//CheckUser checks if a user is whitelisted
func (a *DBAgent) CheckUser(ctx context.Context, uuid string) (bool, error) {
	var tmp interface{}
	err := a.db.QueryRowContext(ctx, `
SELECT true FROM users
WHERE "uuid" = $1`,
		uuid,
	).Scan(&tmp)
	fmt.Println("checking user", uuid, tmp)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "failed to check user authorization")
	}
	return true, nil
}
