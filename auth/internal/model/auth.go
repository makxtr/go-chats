package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/note_v1.NoteV1/Get"
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
