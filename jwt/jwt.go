package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dokidokikoi/go-common/errors"
)

type CustomClaims struct {
	ID    uint   `json:"id"`
	Emial string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(claims CustomClaims, signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func VerifyToken(tokenString string, signingKey string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.ErrPlzLogin
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				return nil, errors.ErrPlzLogin
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, errors.ErrPlzLogin
			} else {
				return nil, ve
			}
		}
		return nil, errors.ErrPlzLogin
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrPlzLogin
}
