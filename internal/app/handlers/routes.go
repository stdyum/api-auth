package handlers

import (
	"github.com/stdyum/api-common/hc"
	"github.com/stdyum/api-common/http/middlewares"
	"github.com/stdyum/api-common/proto/impl/auth"
	"google.golang.org/grpc"
)

func (h *http) ConfigureRoutes() *hc.Engine {
	engine := hc.New()

	group := engine.Group("api/v1")
	group.Use(hc.Logger(), hc.Recovery(), middlewares.ErrorMiddleware())

	group.GET("auth/oauth2/:provider", h.AuthViaOAuth2)
	group.GET("auth/oauth2/:provider/callback", h.AuthViaOAuth2Callback)

	group.POST("signup", h.SignUp)
	group.POST("login", h.Login)
	group.GET("self", h.GetSelfUser)

	group.POST("token/update", h.UpdateToken)

	group.POST("confirm/email/code", h.ConfirmEmailByCode)

	group.POST("reset/password/request", h.ResetPasswordRequest)
	group.POST("reset/password/code", h.ResetPasswordByCode)

	return engine
}

func (h *gRPC) ConfigureRoutes() *grpc.Server {
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServer(grpcServer, h)
	return grpcServer
}
