package user

import (
	"auth/internal/converter"
	desc "auth/pkg/user_v1"
	"context"
	"log"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, mapError(err)
	}

	log.Printf("id: %d, name: %s, email: %s, created_at: %v, updated_at: %v\n", user.ID, user.Info.Name, user.Info.Email, user.CreatedAt, user.UpdatedAt)

	return &desc.GetResponse{
		User: converter.ToUserFromService(user),
	}, nil
}
