package controllers

import (
	"kanban-app/api/models"
	"kanban-app/api/services"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var userService *services.UserService

func init() {
	userService = services.NewUserService()
}

// RegisterUser handles user registration.
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration details"
// @Success 201 {object} models.User "User registered successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := userService.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to register user: " + err.Error()})
		return
	}

	user.Password = "" // TODO check if password wont show by default due to json binding in struct
	c.JSON(http.StatusCreated, user)
}

// LoginUser handles user login and generates a JWT.
// @Summary Log in a user
// @Description Authenticate user and return a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User login credentials"
// @Success 200 {object} models.TokenResponse "Successfully logged in"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /login [post]
func LoginUser(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := userService.AuthenticateUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: err.Error()})
		return
	}

	token, err := GenerateJwt(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{Token: token})
}

func GenerateJwt(userID string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
