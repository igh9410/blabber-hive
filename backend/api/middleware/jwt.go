// middleware/jwt.go

package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Aud         string `json:"aud"`
	Exp         int64  `json:"exp"`
	Iat         int64  `json:"iat"`
	Iss         string `json:"iss"`
	Sub         string `json:"sub"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	Aal         string `json:"aal"`
	SessionID   string `json:"session_id"`
	AppMetadata struct {
		Provider  string   `json:"provider"`
		Providers []string `json:"providers"`
	} `json:"app_metadata"`
	UserMetadata struct{} `json:"user_metadata"`
	Amr          []struct {
		Method    string `json:"method"`
		Timestamp int64  `json:"timestamp"`
	} `json:"amr"`
	jwt.RegisteredClaims
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func EnsureValidToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization") // Get token from the Authorization header
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// If the token string starts with "Bearer ", remove it.
		if len(tokenString) > 7 && strings.ToUpper(tokenString[0:7]) == "BEARER " {
			tokenString = tokenString[7:]
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("AllYourBase"), nil
		})

		if token.Valid {
			fmt.Println("You look nice today")
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			fmt.Println("That's not even a token")
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			// Invalid signature
			fmt.Println("Invalid signature")
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	}

}
