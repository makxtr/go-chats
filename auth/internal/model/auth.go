package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/chat_server_v1.ChatServerV1/SendMessage"
)

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginUserCommand struct {
	Username string
	Password string
}

type RefreshToken struct {
	Token string
}
