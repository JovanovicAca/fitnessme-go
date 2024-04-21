package utils

import (
	"errors"
	"fitnessme/usermanagement/pkg/models"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type JwtClaims struct {
	jwt.StandardClaims
	Id    uuid.UUID
	Email string
	Role  string
}

var jwtKey = []byte("my_secret_key")

func HashPassword(pass string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (jw *JwtWrapper) GenerateToken(user models.User) (signedToken string, err error) {
	claims := &JwtClaims{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(jw.ExpirationHours)).Unix(),
			Issuer:    jw.Issuer,
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = t.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (jw *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*JwtClaims)

	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	return claims, nil
}
