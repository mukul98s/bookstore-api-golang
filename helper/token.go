package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetTokens(id string) (string, error) {
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := tokenClaim.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return string(token), err
}

func GetRefreshToken(id string) (string, error) {
	refreshTokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshToken, err := refreshTokenClaim.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return string(refreshToken), err
}
