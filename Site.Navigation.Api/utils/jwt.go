package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte
var jwtExpireHours = 24

type Claims struct {
	UserID  int    `json:"user_id"`
	Account string `json:"account"`
	jwt.RegisteredClaims
}

func InitJWT(secret string, expireHours int) {
	jwtSecret = []byte(secret)
	if expireHours > 0 {
		jwtExpireHours = expireHours
	}
}

func GenerateToken(userID int, account string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(jwtExpireHours) * time.Hour)

	claims := Claims{
		UserID:  userID,
		Account: account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			Issuer:    "site-navigation-api",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtSecret)
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, errors.New("无效的token")
}
