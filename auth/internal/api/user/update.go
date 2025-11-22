package user

import (
	"auth/internal/converter"
	desc "auth/pkg/user_v1"
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, req.GetId(), converter.ToUserUpdateFromDesc(req))
	if err != nil {
		return nil, mapError(err)
	}

	log.Printf("updated user with id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}
