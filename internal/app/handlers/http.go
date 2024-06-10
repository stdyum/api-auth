package handlers

import (
	"github.com/stdyum/api-auth/internal/app/controllers"
	"github.com/stdyum/api-common/hc"
	confHttp "github.com/stdyum/api-common/http"
)

type HTTP interface {
	confHttp.Routes

	AuthViaOAuth2(ctx *hc.Context)
	AuthViaOAuth2Callback(ctx *hc.Context)

	SignUp(ctx *hc.Context)
	Login(ctx *hc.Context)
	GetSelfUser(ctx *hc.Context)

	UpdateToken(ctx *hc.Context)

	ConfirmEmailByCode(ctx *hc.Context)

	ResetPasswordRequest(ctx *hc.Context)
	ResetPasswordByCode(ctx *hc.Context)
}

type http struct {
	controller  controllers.Controller
	redirectURL string
}

func NewHTTP(controller controllers.Controller, redirectURL string) HTTP {
	return &http{
		controller:  controller,
		redirectURL: redirectURL,
	}
}
