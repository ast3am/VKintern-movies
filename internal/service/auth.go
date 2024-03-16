package service

import (
	"context"
	"errors"
	"github.com/ast3am/VKintern-movies/internal/utils"
)

func (s *Service) Auth(ctx context.Context, email, password string) (token string, err error) {
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	credsCorrect := user.CheckCreds(email, password)

	if !credsCorrect {
		return "", errors.New("wrong email or password")
	}

	token, err = utils.GetToken(user.Email, user.Role)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) CheckToken(token, permissionLevel string) error {
	err := utils.CheckPermissionByToken(token, permissionLevel)
	return err
}
