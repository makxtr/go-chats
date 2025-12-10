package access

import (
	"auth/internal/model"
	"auth/internal/utils"
	"context"
)

var accessibleRoles map[string]string

func (s *serv) Check(ctx context.Context, accessToken, endPoint string) (bool, error) {
	claims, err := utils.VerifyToken(accessToken, []byte(s.securityConfig.AccessKey()))
	if err != nil {
		return false, TokenInvalid
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return false, RolesGetFailed
	}

	role, ok := accessibleMap[endPoint]
	if !ok {
		return true, nil
	}

	if role == claims.Role {
		return true, nil
	}

	return false, Denied
}

func (s *serv) accessibleRoles(ctx context.Context) (map[string]string, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)

		accessibleRoles[model.ExamplePath] = "ADMIN"
	}

	return accessibleRoles, nil
}
