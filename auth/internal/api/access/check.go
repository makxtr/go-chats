package access

import (
	desc "auth/pkg/access_v1"
	"context"
	"log"

	"auth/internal/utils"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	accessToken, err := utils.GetTokenFromCtx(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	log.Printf("check access")

	_, err = i.accessService.Check(ctx, accessToken, req.GetEndpointAddress())
	if err != nil {
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
