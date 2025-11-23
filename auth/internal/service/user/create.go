package user

import (
	"auth/internal/model"
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) Create(ctx context.Context, command *model.CreateUserCommand) (int64, error) {
	var id int64

	if err := command.Validate(); err != nil {
		return 0, status.Error(codes.InvalidArgument, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(command.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return 0, status.Error(codes.Internal, "failed to hash password")
	}

	createData := &model.CreateUserData{
		Info:           command.Info,
		HashedPassword: string(hashedPassword),
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.Create(ctx, createData)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.userRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &model.UserLog{
			Action:   "user_created",
			EntityID: id,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
