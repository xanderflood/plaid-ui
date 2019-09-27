package auth

import (
	"errors"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//AuthorizationContextKey is the key used to store the Authorization
//in the context
const AuthorizationContextKey = "PLAID_UI_PUBLIC_API_AUTHORIZATION"

//Authorization describes the authorities stored in a user JWT
type Authorization struct {
	jwt.StandardClaims
	UserUUID *string `json:"sub,omitempty"`
	IsAdmin  bool    `json:"pld.admn"`
}

//Valid applies standard JWT validations as well as generic
//user authorization rules.
func (a *Authorization) Valid() error {
	if err := a.StandardClaims.Valid(); err != nil {
		return err
	}

	if a.UserUUID != nil && len(*a.UserUUID) != 0 && !a.IsAdmin {
		return errors.New("non-admin user does not have an identity")
	}

	return nil
}

//AuthGetter is a helper for grabbing the Authorization
//that the middleware stores in the context.
//go:generate counterfeiter . AuthGetter
type AuthGetter func(c *gin.Context) (Authorization, bool)

//GetAuthorization is the default AuthGetter
func GetAuthorization(c *gin.Context) (Authorization, bool) {
	authIface := c.Value(AuthorizationContextKey)
	auth, ok := authIface.(Authorization)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization object was found"})
	}
	return auth, ok
}

//Authorizer represents the needed interactions with jwt.Parser
//go:generate counterfeiter . Authorizer
type Authorizer interface {
	ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
}

//JWTMiddleware checks for a JWT in a bearer token on the request
//and converts it into an Authorzation struct, which is stored in
//the context.
func JWTMiddleware(signingSecret string, authorizer Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authorization provided"})
			return
		}

		var auth Authorization
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := authorizer.ParseWithClaims(tokenString, &auth, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("must use HMAC signing")
			}

			return signingSecret, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set(AuthorizationContextKey, auth)
		c.Next()
	}
}
