package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//RegistrationRequest encodes a single request for user registration
type RegistrationRequest struct {
	UserUUID string `json:"user_uuid" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

//RegisterUser adds all the accounts associated with this plaid item
func (a ServerAgent) RegisterUser(c *gin.Context) {
	authorization, ok := a.authorize(c)
	if !ok {
		return
	}

	if !authorization.Admin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "this endpoint is only accessible to users with administrative priveleges"})
		return
	}

	var req RegistrationRequest
	err := c.Bind(&req)
	if err != nil {
		//TODO is this error message safe to expose?
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = a.dbClient.RegisterUser(c, req.UserUUID, req.Email)
	if err != nil {
		a.logger.Errorf("register user `%s` failed: %s", req.UserUUID, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user registration failed - see logs for details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_uuid": req.UserUUID,
		"email":     req.Email,
	})
}
