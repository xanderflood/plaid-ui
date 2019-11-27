package auth

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/xanderflood/plaid-ui/cmd/api/server/views"
	"github.com/xanderflood/plaid-ui/lib/tools"
	"github.com/xanderflood/plaid-ui/pkg/db"
)

//AuthorizationContextKey is the key used to store the Authorization
//in the context
const AuthorizationContextKey = "PLAID_UI_PUBLIC_API_AUTHORIZATION"

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

//Getter is a helper for grabbing the Authorization
//that the middleware stores in the context.
//go:generate counterfeiter . Getter
type Getter func(c *gin.Context) (Authorization, bool)

//GetAuthorizationFromContext is the default Getter
func GetAuthorizationFromContext(c *gin.Context) (Authorization, bool) {
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

//AuthorizationManager exposes middleware functionality for authorization
type AuthorizationManager interface {
	BackendMiddleware() gin.HandlerFunc
	FrontendMiddleware() gin.HandlerFunc
}

//JWTAuthorizationManager provides a JWT-based implementation of AuthorizationManager
type JWTAuthorizationManager struct {
	logger          tools.Logger
	renderer        views.Renderer
	signingSecret   string
	authorizer      Authorizer
	db              db.DB
	loginBaseURLRef *url.URL
}

//NewAuthorizationManager creates a new JWTAuthorizationManager
func NewAuthorizationManager(
	logger tools.Logger,
	renderer views.Renderer,
	signingSecret string,
	authorizer Authorizer,
	db db.DB,
	loginBaseURLRef *url.URL,
) JWTAuthorizationManager {
	return JWTAuthorizationManager{
		logger:          logger,
		renderer:        renderer,
		signingSecret:   signingSecret,
		authorizer:      authorizer,
		db:              db,
		loginBaseURLRef: loginBaseURLRef,
	}
}

func (a JWTAuthorizationManager) getAuthorizationFromString(tokenString string) (Authorization, error) {
	var auth Authorization
	_, err := a.authorizer.ParseWithClaims(tokenString, &auth, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("must use HMAC signing")
		}

		return []byte(a.signingSecret), nil
	})
	if err != nil {
		return Authorization{}, err
	}

	return auth, nil
}

func (a JWTAuthorizationManager) requireServiceAccess(c *gin.Context, auth Authorization) (bool, error) {
	if auth.User || auth.Admin {
		return true, nil
	}

	return a.db.CheckUser(c, auth.UserUUID)
}

//BackendMiddleware checks for a JWT in a bearer token on the request
//and converts it into an Authorzation struct, which is stored in
//the context.
func (a JWTAuthorizationManager) BackendMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authorization provided"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		auth, err := a.getAuthorizationFromString(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ok, err := a.requireServiceAccess(c, auth)
		if err != nil {
			a.logger.Errorf(err.Error())
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": "An internal error occurred"},
			)
			return
		}
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "The administrator has not granted you access to this service."},
			)
			return
		}

		c.Set(AuthorizationContextKey, auth)
		c.Next()
	}
}

//FrontendMiddleware checks for a JWT in a request token. If it's not
//there or invalid, redirect the user to the login flow, with instructions
//to refer the user back here afterwards.
func (a JWTAuthorizationManager) FrontendMiddleware() gin.HandlerFunc {
	redirectToLogin := func(c *gin.Context) {
		//build the URL to redirect back to after logging ing
		requestURL := (&url.URL{
			Scheme: "https",
			Host:   c.Request.Host,
			Path:   c.Request.URL.Path,
		})

		//prepare it for the query string
		query := make(url.Values)
		query.Set("referrer_url", requestURL.String())

		//copy the base URL and add the query param
		loginBaseURLObj := *a.loginBaseURLRef
		loginBaseURL := &loginBaseURLObj
		loginBaseURL.RawQuery = query.Encode()

		c.Redirect(http.StatusTemporaryRedirect, loginBaseURL.String())
		c.Abort()
	}

	return func(c *gin.Context) {
		jwtCookie, err := c.Request.Cookie("_identify_jwt_string")
		if err != nil {
			redirectToLogin(c)
			return
		}

		auth, err := a.getAuthorizationFromString(jwtCookie.Value)
		if err != nil {
			redirectToLogin(c)
			return
		}

		ok, err := a.requireServiceAccess(c, auth)
		if err != nil {
			a.logger.Errorf(err.Error())
			a.renderer.RenderStatusCode(
				http.StatusInternalServerError,
				"An internal error occurred",
				c,
			)
			return
		}
		if !ok {
			a.renderer.RenderNotRegistered(auth.Email, c)
			return
		}

		c.Next()
	}
}
