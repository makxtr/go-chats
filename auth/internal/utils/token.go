package utils

import (
	"auth/internal/model"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

const authPrefix = "Bearer "

func GenerateToken(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: info.Name,
		Role:     info.Role.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	return claims, nil
}

func GetTokenFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrNoMetadata
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", ErrNoHeader
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return "", fmt.Errorf("%w: header value was '%s'", ErrInvalidFormat, authHeader[0])
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	return accessToken, nil
}
