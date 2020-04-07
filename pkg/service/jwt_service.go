package service

import (
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/sarulabs/di"
)

// JwtError contains error information for JWT token validation
type JwtError struct {
	// HTTP status code (usually 401 or 403)
	StatusCode int
	// Error message
	Message string
}

// NewJwtError creates new instance of JwtError
func NewJwtError(statusCode int, message string) *JwtError {
	return &JwtError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// Jwt provides methods to create and to parse JWT tokens
type Jwt interface {
	// Validate a JWT token
	ValidateToken(token string) (*model.User, *JwtError)

	// Generate new JWT token
	GenerateToken(user *model.User) (string, error)
}

const jwtKey = "JwtService"

// GetJwt returns an implementation of Jwt from DI container
func GetJwt(c di.Container) Jwt {
	return c.Get(jwtKey).(Jwt)
}

type jwtService struct {
	logger         *log.Logger
	tokenPassword  string
	signingMethod  jwt.SigningMethod
	userRepository UserRepository
}

// Validate a JWT token
func (s *jwtService) ValidateToken(tokenStr string) (*model.User, *JwtError) {
	if tokenStr == "" {
		return nil, NewJwtError(401, "missing access token")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.tokenPassword), nil
	})

	if err != nil || !token.Valid || token.Claims.Valid() != nil {
		return nil, NewJwtError(401, "bad access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		id, exists := claims["id"]
		if exists {
			idStr, ok := id.(string)
			if ok {
				user, err := s.userRepository.Get(idStr)
				if err == nil {
					return user, nil
				}
			}
		}
	}

	return nil, NewJwtError(401, "bad access token")
}

// Generate new JWT token
func (s *jwtService) GenerateToken(user *model.User) (string, error) {
	token := jwt.New(s.signingMethod)
	token.Claims = jwt.MapClaims{
		"id": user.UserName,
	}
	str, err := token.SignedString([]byte(s.tokenPassword))
	if err != nil {
		s.logger.Printf("unable to generate jwt token for user #%d: %v", user.ID, err)
	}
	return str, err
}
