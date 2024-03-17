package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var roleMap = map[string]int{
	"user":  0,
	"admin": 1,
}

func GetToken(email, role string) (string, error) {
	secretKey := []byte("my_secret_key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	resToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("can't make token")
	}
	return resToken, nil
}

func CheckPermissionByToken(token, permissionLevel string) error {
	secretKey := []byte("my_secret_key")
	tokenParse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method from token")
		}
		return secretKey, nil
	})

	if err != nil {
		err = errors.New("not a valid token")
		return err
	}
	role := ""
	if claims, ok := tokenParse.Claims.(jwt.MapClaims); ok && tokenParse.Valid {
		role = claims["role"].(string)
	} else {
		err = errors.New("not a valid token")
		return err
	}
	if roleMap[role] < roleMap[permissionLevel] {
		err = errors.New("permission denied")
		return err
	}
	return nil
}
