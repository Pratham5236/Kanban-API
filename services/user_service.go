package services

import (
	"errors"
	"fmt"
	"kanban-app/api/database"
	"kanban-app/api/models"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) RegisterUser(req models.RegisterRequest) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user %w", result.Error)
	}

	log.Printf("User registered: %s (%s)\n", user.Username, user.Email)
	return &user, nil
}

func (s *UserService) AuthenticateUser(req models.LoginRequest) (*models.User, error) {
	var user models.User

	result := database.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("database error during authentication: %w", result.Error)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("password comparison failed: %w", err)
	}

	return &user, nil
}
