package helpers

import (
	"kanban-task/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Username  string    `json:"username"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	jwt.StandardClaims
}
type OutputToken struct {
	Username  string    `json:"username"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func GenerateToken(user model.User) (*OutputToken, error) {

	expireTime := time.Now().Add(time.Hour * 48)
	claims := &Claims{
		Username:  user.Username,
		UserID:    user.UserID,
		CreatedAt: user.CreatedAt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	outputToken := &OutputToken{
		Username:  user.Username,
		UserID:    user.UserID,
		CreatedAt: user.CreatedAt,
		Token:     tokenString,
		ExpiresAt: expireTime,
	}
	return outputToken, err

}

func ValidateToken(tokenValue string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Minute {
		return token, claims, err
	}
	expireTime := time.Now().Add(2 * time.Hour)
	claims.ExpiresAt = expireTime.Unix()
	tkn := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tkn.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tkn, claims, err

}
