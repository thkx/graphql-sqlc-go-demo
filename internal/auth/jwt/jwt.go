package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey           []byte
	accessTokenDuration time.Duration
)

func SetupJWTConfig(secret []byte, duration time.Duration) {
	secretKey = secret
	accessTokenDuration = duration
}

// GenerateToken 生成 JWT access token
func GenerateToken(userID, name, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    name,
		"email":   email,
		"exp":     time.Now().Add(accessTokenDuration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateToken 解析并验证 JWT，返回用户ID（或完整 claims）
func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return "", errors.New("user_id missing in token")
	}

	return userID, nil
}

// Optional: GetClaims 返回完整 claims（用于需要 name/email 的场景）
func GetClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}
