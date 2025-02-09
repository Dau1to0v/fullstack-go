package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/Dau1to0v/fullstack-go/pkg/repository"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	salt       = "fhaskh124khkhf9"
	signingKey = "afsf13ewfr23fwef"
	tokenTTL   = time.Hour * 12
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePassword(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetById(id int) (models.User, error) {
	return s.repo.GetUserById(id)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePassword(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		logrus.Errorf("Ошибка при парсинге токена: %v", err)
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		logrus.Error("Клеймы токена неверного типа или токен недействителен")
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	logrus.Infof("Извлечён userId из токена: %d", claims.UserId)
	return claims.UserId, nil
}

func generatePassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
