package handlers

import (
	"github.com/stdyum/api-auth/internal/app/controllers"
	"github.com/stdyum/api-common/grpc"
	"github.com/stdyum/api-common/proto/impl/auth"
)

type GRPC interface {
	grpc.Routes
	auth.AuthServer
}

type gRPC struct {
	auth.UnimplementedAuthServer

	controller controllers.Controller
}

func NewGRPC(controller controllers.Controller) GRPC {
	return &gRPC{
		controller: controller,
	}
}
