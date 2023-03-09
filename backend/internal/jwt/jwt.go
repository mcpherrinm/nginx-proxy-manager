package jwt

import (
	"time"

	"npm/internal/entity/user"
	"npm/internal/logger"

	"github.com/dgrijalva/jwt-go"
	"github.com/rotisserie/eris"
)

// UserJWTClaims is the structure of a JWT for a User
type UserJWTClaims struct {
	UserID int      `json:"uid"`
	Roles  []string `json:"roles"`
	jwt.StandardClaims
}

// GeneratedResponse is the response of a generated token, usually used in http response
type GeneratedResponse struct {
	Expires int64  `json:"expires"`
	Token   string `json:"token"`
}

// Generate will create a JWT
func Generate(userObj *user.Model, forSSE bool) (GeneratedResponse, error) {
	var response GeneratedResponse

	key, _ := GetPrivateKey()
	expires := time.Now().AddDate(0, 0, 1) // 1 day
	issuer := "api"

	if forSSE {
		issuer = "sse"
	}

	// Create the Claims
	claims := UserJWTClaims{
		userObj.ID,
		[]string{"user"},
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expires.Unix(),
			Issuer:    issuer,
		},
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	var err error
	token.Signature, err = token.SignedString(key)
	if err != nil {
		logger.Error("JWTError", eris.Wrapf(err, "Error signing token: %v", err))
		return response, err
	}

	response = GeneratedResponse{
		Expires: expires.Unix(),
		Token:   token.Signature,
	}

	return response, nil
}