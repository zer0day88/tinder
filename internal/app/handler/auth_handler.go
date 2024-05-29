package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zer0day88/tinder/helper"
	"github.com/zer0day88/tinder/internal/app/model"
	"github.com/zer0day88/tinder/internal/app/service"
	"github.com/zer0day88/tinder/pkg/response"
)

type authHandler struct {
	authSrv *service.AuthService
	log     zerolog.Logger
}

type AuthHandler interface {
	SignUp(c echo.Context) error
	Login(c echo.Context) error
	Cek(c echo.Context) error
}

func NewAuthHandler(log zerolog.Logger, authSrv *service.AuthService) AuthHandler {
	return &authHandler{log: log, authSrv: authSrv}
}

func (h *authHandler) SignUp(c echo.Context) error {
	var req model.SignUpRequest

	if err := c.Bind(&req); err != nil {
		return response.WriteJSON(c, response.ErrBadRequest)
	}

	errv := helper.ValidateStruct(req)

	if errv != nil {
		return response.WriteJSON(c, response.ErrBadRequest.WithMsg(*errv))
	}

	resp := h.authSrv.SignUp(c.Request().Context(), req)

	if resp.Err != nil {
		h.log.Err(resp.Err).Send()
	}

	return response.WriteJSON(c, resp)
}

func (h *authHandler) Login(c echo.Context) error {
	var req model.SigningRequest

	if err := c.Bind(&req); err != nil {
		return response.WriteJSON(c, response.ErrBadRequest)
	}

	resp, api := h.authSrv.Login(c.Request().Context(), req)

	if api.Err != nil {
		h.log.Err(api.Err).Send()
	}

	return response.WriteJSON(c, api.WithData(resp))
}

func (h *authHandler) Cek(c echo.Context) error {

	return response.WriteJSON(c, response.OKNoError)
}
