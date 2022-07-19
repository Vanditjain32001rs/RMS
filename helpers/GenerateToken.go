package helpers

import (
	"RMS/models"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

var SecretKey = "my_secret_key"

func GetSecretKey() string {
	return SecretKey
}

var ExpirationTime = time.Now().Add(time.Hour * 24 * 7)

func TokenGeneration(username string, user *models.UserRoleID) (string, error) {

	claims := &models.Claims{
		Username: username,
		UserRole: user.UserRole,
		UserID:   user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: ExpirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		log.Printf("TokenGeneration : Error in signing the token.")
		return "", err
	}
	return signedToken, nil
}
