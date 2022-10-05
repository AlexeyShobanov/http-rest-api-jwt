package pkg

import (
	"errors"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type ConfToken struct {
	SigningKey string `toml:"signing_key"`
	TokenTTL   int    `toml:"token_ttl"`
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int
}

func Init(pathConfig string) *ConfToken {
	t := &ConfToken{}
	_, err := toml.DecodeFile(pathConfig, t)
	if err != nil {
		logrus.Fatal(err)
	}
	return t
}

func (t *ConfToken) GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(t.TokenTTL) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})

	return token.SignedString([]byte(t.SigningKey))
}

func (t *ConfToken) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(t.SigningKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}
