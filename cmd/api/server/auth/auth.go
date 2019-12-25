package auth

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/xanderflood/plaid-ui/pkg/db"
)

//Authorization describes the authorities stored in a user JWT
type Authorization struct {
	jwt.StandardClaims
	UserUUID string `json:"sub,omitempty"`
	Email    string `json:"eml,omitempty"`
	Admin    bool   `json:"login.adm,omitempty"`
	User     bool   `json:"login.user,omitempty"`
}

//Valid applies standard JWT validations as well as generic
//user authorization rules.
func (a *Authorization) Valid() error {
	if err := a.StandardClaims.Valid(); err != nil {
		return err
	}

	if !a.Admin && len(a.UserUUID) == 0 {
		return errors.New("non-admin user does not have an identity")
	}

	return nil
}

func (a Authorization) GetDBAuthorization(userUUID string) (db.Authorization, error) {
	if a.Admin || a.UserUUID == userUUID {
		return db.Authorization{UserUUID: userUUID}, nil
	}
	return db.Authorization{}, fmt.Errorf("refusing to produce database authorization for user UUID `%s`", userUUID)
}
