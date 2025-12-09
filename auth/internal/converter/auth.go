package converter

import (
	"auth/internal/model"
	desc "auth/pkg/auth_v1"
)

func ToLoginFromDesc(req *desc.LoginRequest) *model.LoginUserCommand {
	return &model.LoginUserCommand{
		Username: req.Username,
		Password: req.Password,
	}
}
