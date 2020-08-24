package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
	"vehicledb/db"
)

var secretToken = []byte("sup3rs3cr3t")

func CreateToken(user *db.User) (string, error) {
	claims := ClaimsUser{
		EmailAddress: user.EmailAddress,
		UserID:       user.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "vehicledb",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretToken)
	return signedToken, err
}

func ValidateToken(token string) (*ClaimsUser, error) {
	claims := &ClaimsUser{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretToken, nil
	})

	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil
}
