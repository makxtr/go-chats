package user

import (
	desc "auth/pkg/user_v1"
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, mapError(err)
	}

	log.Printf("deleted user with id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}
