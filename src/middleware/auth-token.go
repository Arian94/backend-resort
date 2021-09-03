package middleware

import (
	signupLogin "Resort/src/signup-login"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	BEARER_SCHEMA = "Bearer"
)

func AuthorizeJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "No Auth header provided")
	}
	tokenString := authHeader[len(BEARER_SCHEMA)+1:]
	if tokenString == "null" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Token is null")
	}

	if token, err := signupLogin.JWTAuthService().ValidateToken(tokenString); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
	} else {
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["email"])
	}
}

func AuthorizeOptionalJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			AuthorizeJWT(c)
		}
	}
}
