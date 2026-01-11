package auth

import (
	"subscribly/customerrors"
	"subscribly/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const SECRET_KEY = "alksjd;faj9"

func GenerateAccesndRefresh(userId int, role string) (string, string, error) {

	newClaims := models.Claims{
		UserId: userId,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "subscribly",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	accessToken, err := token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(userId, role)
	if err != nil{
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func GenerateRefreshToken(userId int, role string) (string, error) {
	
	newClaims := models.Claims{
		UserId: userId,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "subscribly",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute *15)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	signedToken, err := token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", err
	}

	return signedToken, nil

}


func ValidateJwtToken(mytoken string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(mytoken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return claims, err
	}

	if !token.Valid {
		return claims, customerrors.ErrInvalidToken
	}
	return claims, nil
}

func ValidateRefreshToken(refreshToken string) (error) {

	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil{
		return err
	}

	if !token.Valid{
		return customerrors.ErrInvalidToken
	}

	return nil
}



func GenerateNewToken(userId int, role string) (string, error){
	claims := &models.Claims{
		UserId: userId,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "subscribly",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken,err := token.SignedString([]byte(SECRET_KEY))
	if err != nil{
		return "", err
	}

	return signedToken, nil

}