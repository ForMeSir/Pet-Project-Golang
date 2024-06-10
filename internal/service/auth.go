package service

import (
	"errors"
	"pet/internal/model"
	"pet/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	AcTimeLive  = 20 * time.Minute
	RefTimeLive = 24 * time.Hour
	signingKey  = "ghhfjjri38"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId    uuid.UUID `json:"user_id"`
	SessionId uuid.UUID `json:"session_id"`
	UserRole  string    `json:"user_role"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (a *AuthService) CreateUser(name string, username string, password string) (id uuid.UUID, err error) {
	userRole := "user"
	return a.repo.CreateUser(name, username, password, userRole)
}

func (a *AuthService) CreateAdmin(name string, username string, password string) (id uuid.UUID, err error) {
	userRole := "admin"
	return a.repo.CreateUser(name, username, password, userRole)
}

func (a *AuthService) GetUser(username string, password string) (user model.User, err error) {
	return a.repo.GetUser(username, password)
}

func (a *AuthService) GenerateToken(userId uuid.UUID, userRole string) (refreshtoken string, accesstoken string, err error) {
	id := uuid.New()
	actoken := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AcTimeLive).Unix(),
			IssuedAt:  time.Now().Unix(),
		}, userId, id, userRole,
	})
	accesstoken, err = actoken.SignedString([]byte(signingKey))
	if err != nil {
		return
	}

	reftoken := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RefTimeLive).Unix(),
			IssuedAt:  time.Now().Unix(),
		}, userId, id, userRole,
	})

	refreshtoken, err = reftoken.SignedString([]byte(signingKey))
	if err != nil {
		return
	}

	err = a.repo.CreateSession(userId, userRole, id, refreshtoken)
	if err != nil {
		return
	}
	return
}

func (a *AuthService) ParseToken(token string) (ptoken model.Token, err error) {
	parsedtoken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(parsedtoken *jwt.Token) (interface{}, error) {
		if _, ok := parsedtoken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return
	}

	claims, ok := parsedtoken.Claims.(*tokenClaims)
	if !ok {
		return ptoken, errors.New("token claims are not of type *tokenClaims")
	}

	_, err = a.repo.FindSession(claims.SessionId)
	if err != nil {
		return
	}
	ptoken.UserID = claims.UserId
	ptoken.SessionID = claims.SessionId
	ptoken.UserRole = claims.UserRole
	ptoken.CreatedAt = time.Unix(claims.IssuedAt, 0)
	ptoken.ExpiratedAt = time.Unix(claims.ExpiresAt, 0)
	return
}

func (a *AuthService) Refresh(refreshtoken string, accesstoken string) (actoken string, err error) {
	parsedtoken, err := jwt.ParseWithClaims(refreshtoken, &tokenClaims{}, func(parsedtoken *jwt.Token) (interface{}, error) {
		if _, ok := parsedtoken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return
	}
	claims, ok := parsedtoken.Claims.(*tokenClaims)
	if !ok {
		return accesstoken, errors.New("token claims are not of type *tokenClaims")
	}

	session, err := a.repo.FindSession(claims.SessionId)

	if err != nil {
		return
	}

	if session.Refreshtoken != refreshtoken {
		return actoken, errors.New("invalid token")
	}

	if session.IsBlocked {
		return actoken, errors.New("token blocked")
	}

	acstoken := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AcTimeLive).Unix(),
			IssuedAt:  time.Now().Unix(),
		}, claims.UserId, claims.SessionId, claims.UserRole,
	})
	actoken, err = acstoken.SignedString([]byte(signingKey))
	if err != nil {
		return
	}

	return
}
