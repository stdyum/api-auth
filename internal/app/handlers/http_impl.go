package handlers

import (
	"github.com/stdyum/api-auth/internal/app/dto"
	"github.com/stdyum/api-common/hc"

	netHttp "net/http"
)

func (h *http) SignUp(ctx *hc.Context) {
	var signUpDTO dto.SignUpRequestDTO
	if err := ctx.BindJSON(&signUpDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	responseDTO, err := h.controller.SignUp(ctx, signUpDTO)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(netHttp.StatusCreated, responseDTO)
}

func (h *http) Login(ctx *hc.Context) {
	var loginDTO dto.LoginRequestDTO
	if err := ctx.BindJSON(&loginDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	responseDTO, err := h.controller.Login(ctx, loginDTO)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(netHttp.StatusOK, responseDTO)
}

func (h *http) UpdateToken(ctx *hc.Context) {
	var requestDTO dto.UpdateTokenRequestDTO
	if err := ctx.BindJSON(&requestDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	responseDTO, err := h.controller.UpdateToken(ctx, requestDTO)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(netHttp.StatusOK, responseDTO)
}

func (h *http) ConfirmEmailByCode(ctx *hc.Context) {
	var codeDTO dto.ConfirmEmailByCodeRequestDTO
	if err := ctx.BindJSON(&codeDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := h.controller.ConfirmEmailByCode(ctx, codeDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(netHttp.StatusNoContent)
}

func (h *http) ResetPasswordRequest(ctx *hc.Context) {
	var passwordRequestDTO dto.ResetPasswordRequestDTO
	if err := ctx.BindJSON(&passwordRequestDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := h.controller.ResetPasswordRequest(ctx, passwordRequestDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(netHttp.StatusNoContent)
}

func (h *http) ResetPasswordByCode(ctx *hc.Context) {
	var passwordDTO dto.ResetPasswordByCodeRequestDTO
	if err := ctx.BindJSON(&passwordDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := h.controller.ResetPasswordByCode(ctx, passwordDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(netHttp.StatusNoContent)
}
