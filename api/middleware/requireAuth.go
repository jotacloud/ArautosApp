package middleware

import (
	"ArautosApp/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"ArautosApp/initializers"
	"ArautosApp/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": utils.ErrTokenNotFound})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %v", utils.ErrUnexpectedSigning, token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err) // Logging the error
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidToken})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": utils.ErrExpiredToken})
			return
		}

		var user models.User
		if err := initializers.DB.First(&user, claims["sub"]).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": utils.ErrUserNotFound})
			return
		}

		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidToken})
	}
}
