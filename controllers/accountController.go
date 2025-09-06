package controllers

import (
	"log"
	"net/http"
	"os"
	"serverless/initializers"
	"serverless/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountBody struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type EmailChange struct {
	Email string `json:"email" binding:"required,email,max=255"`
}

type AccountResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	//Posts     []ResponsePost `json:"posts"`
	UpdatedAt time.Time `json:"updated_at"`
}

func AccountCreate(c *gin.Context) {
	var req AccountBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request Object",
		})
		return
	}

	// Hash Password and normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to Hash Password",
		})
		return
	}

	var existingAccount models.Account
	if err := initializers.DB.Where("email = ?", req.Email).First(&existingAccount).Error; err != nil {
		log.Printf("Account already exists for email: %s", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email Already Exists",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("Error checking existing account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error Checking account",
		})
		return
	}

	account := models.Account{
		Email:    req.Email,
		Password: string(hash),
	}

	if err := initializers.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error Creating account",
		})
		return
	}

	resp := AccountResponse{
		ID:        account.ID,
		Email:     account.Email,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
	c.JSON(http.StatusCreated, gin.H{
		"account": resp,
	})
}

func AccountLogin(c *gin.Context) {
	// Get and sanitize the input request
	var req AccountBody

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("invalid request Object")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request Object",
		})
		return
	}

	// fetch user and compare password hash

	var existingAccount models.Account
	if err := initializers.DB.Where("email = ?", req.Email).First(&existingAccount).Error; err != nil {
		log.Println("Account Does not Exist")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Login",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(existingAccount.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Login",
		})
		return
	}
	// Successfull password hash compared, next is generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": existingAccount.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	log.Printf("Token: %v", token)
	// Sign and get the complete encoded token as string using secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		log.Printf("Unable to sign jwt token with secret %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something Unexpected Occurred",
		})
		return
	}
	// Seed the jwt to Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})

}

func AccountUpdate(c *gin.Context) {
	// Get Auth User from cookie, Auth middleware
	accountID, err := c.Get("accountID")
	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not associated with jwt",
		})
		return
	}

	var req EmailChange
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body"})
		return
	}

	var existingAccount models.Account
	if err := initializers.DB.Where("email = ?", req.Email).First(&existingAccount).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already exists",
		})
	}
	var account models.Account
	if err := initializers.DB.First(&account, accountID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch account",
		})
		return
	}

	if err := initializers.DB.Model(&account).Update("email", req.Email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to Update account",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Account Updated",
	})
}
