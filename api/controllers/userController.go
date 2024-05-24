package controllers

import (
	"ArautosApp/initializers"
	"ArautosApp/models"
	"ArautosApp/utils"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SingUp(c *gin.Context) {

	var body struct {
		Name     string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrFailedToReadBody})
		return
	}

	if !isValidEmail(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidEmail})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrFailedToEncryptPwd})
		return
	}

	user := models.User{Name: body.Name, Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrFailedToCreateUser})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": utils.SuccessUserRegistered})
}

func Login(c *gin.Context) {

	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrFailedToReadBody})

		return
	}

	var user models.User
	initializers.DB.First(&user, "Email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidCredentials})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidPassword})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": utils.ErrFailedToCreateToken,
		})

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
func ValidateLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Você Está Logado!",
	})
}

func isValidEmail(email string) bool {
	// Expressão regular para validar endereços de e-mail
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
