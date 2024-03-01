package handlers

import (
	"context"
	"github.com/stdyum/api-common/grpc"
	"github.com/stdyum/api-common/proto/impl/auth"
)

func (h *gRPC) Auth(ctx context.Context, token *auth.Token) (*auth.User, error) {
	user, err := h.controller.Auth(ctx, token.Token)
	return user, grpc.ConvertError(err)
}
