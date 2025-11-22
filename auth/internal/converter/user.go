package converter

import (
	"auth/internal/model"
	desc "auth/pkg/user_v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserCreateFromDesc(req *desc.CreateRequest) *model.CreateUserCommand {
	return &model.CreateUserCommand{
		Info: model.UserInfo{
			Name:  req.GetName(),
			Email: req.GetEmail(),
			Role:  model.Role(req.GetRole()),
		},
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}
}

func ToUserUpdateFromDesc(req *desc.UpdateRequest) *model.UpdateUserData {
	var name, email *string
	if req.GetInfo().GetName() != nil {
		val := req.GetInfo().GetName().GetValue()
		name = &val
	}
	if req.GetInfo().GetEmail() != nil {
		val := req.GetInfo().GetEmail().GetValue()
		email = &val
	}

	return &model.UpdateUserData{
		Name:  name,
		Email: email,
	}
}

func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Name:      user.Info.Name,
		Email:     user.Info.Email,
		Role:      desc.Role(user.Info.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
