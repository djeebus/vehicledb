package auth

import (
	"github.com/dgrijalva/jwt-go"
	"vehicledb/db"
)

type ClaimsUser struct {
	EmailAddress string
	UserID       db.RowID

	jwt.StandardClaims
}
