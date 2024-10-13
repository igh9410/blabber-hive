// middleware/jwt.go

package middleware

import (
	"log"
	"net/http"
	"os"
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
		var expectedIssuer = os.Getenv("SUPABASE_AUTH")
		var jwtSecretKey = os.Getenv("JWT_SECRET")

		tokenString := extractTokenFromHeader(c)
		if tokenString == "" {
			abortWithUnauthorized(c, "Authorization header required")
			return
		}

		token, err := parseJWT(tokenString, jwtSecretKey)

		if err != nil {

			abortWithUnauthorized(c, err.Error())
			return
		}

		if isValidToken(token, expectedIssuer) {
			email := token.Claims.(jwt.MapClaims)["email"].(string)
			sub := token.Claims.(jwt.MapClaims)["sub"].(string)
			c.Set("email", email)
			c.Set("user_id", sub)
			c.Next()
		} else {
			abortWithUnauthorized(c, "invalid token")
		}

	}
}

// WebSocketAuthMiddleware is a middleware that will check the validity of our JWT for WebSocket endpoints.
func WebSocketAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var expectedIssuer = os.Getenv("SUPABASE_AUTH")
		var jwtSecretKey = os.Getenv("JWT_SECRET")

		tokenString := c.Query("token")
		if tokenString == "" {
			abortWithUnauthorized(c, "token is required")
			return
		}

		token, err := parseJWT(tokenString, jwtSecretKey)

		if err != nil {
			abortWithUnauthorized(c, err.Error())
			return
		}

		if isValidToken(token, expectedIssuer) {
			email := token.Claims.(jwt.MapClaims)["email"].(string)
			sub := token.Claims.(jwt.MapClaims)["sub"].(string)
			c.Set("email", email)
			c.Set("user_id", sub)
			c.Next()
		} else {
			abortWithUnauthorized(c, "invalid token")
		}

	}
}

func extractTokenFromHeader(c *gin.Context) string {
	tokenString := c.GetHeader("Authorization")
	if strings.HasPrefix(strings.ToUpper(tokenString), "BEARER ") {
		return tokenString[7:]
	}
	return tokenString
}

func parseJWT(tokenString, jwtSecretKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
}

func isValidToken(token *jwt.Token, expectedIssuer string) bool {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		issuer := claims["iss"].(string)
		email := claims["email"].(string)
		return issuer == expectedIssuer && email != ""
	}
	return false
}

func abortWithUnauthorized(c *gin.Context, message string) {
	log.Printf("Unauthorized request: %s", message)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
}
