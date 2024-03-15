package handlers

import (
	"context"

	"github.com/stdyum/api-common/grpc"
	"github.com/stdyum/api-common/proto/impl/auth"
)

func (h *gRPC) Auth(ctx context.Context, token *auth.Token) (*auth.User, error) {
	user, err := h.controller.Auth(ctx, token.Token)
	if err != nil {
		return nil, grpc.ConvertError(err)
	}

	return &auth.User{
		Id:            user.ID.String(),
		Login:         user.Login,
		PictureUrl:    user.PictureUrl,
		Email:         user.Email,
		VerifiedEmail: user.VerifiedEmail,
	}, nil
}
