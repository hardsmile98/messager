package handler

import (
	"net/http"

	"gateway/internal/transport/http/dto"
	"gateway/internal/transport/http/helpers"
	"gateway/internal/transport/http/response"
	"gateway/internal/validation"

	pb "github.com/hardsmile98/messager/sdk/auth/v1"
)

type Auth struct {
	client pb.AuthServiceClient
}

func NewAuth(client pb.AuthServiceClient) *Auth {
	return &Auth{client: client}
}

func (h *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.RequestError(w, err)
		return
	}

	resp, err := h.client.Register(r.Context(), &pb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Device:   helpers.DeviceFromRequest(r),
	})
	if err != nil {
		response.GRPCError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"user_id":       resp.UserId,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.RequestError(w, err)
		return
	}

	resp, err := h.client.Login(r.Context(), &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
		Device:   helpers.DeviceFromRequest(r),
	})
	if err != nil {
		response.GRPCError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"user_id":       resp.UserId,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest

	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.RequestError(w, err)
		return
	}

	resp, err := h.client.Logout(r.Context(), &pb.LogoutRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		response.GRPCError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{
		"success": resp.Success,
	})
}

func (h *Auth) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest

	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.RequestError(w, err)
		return
	}

	resp, err := h.client.RefreshToken(r.Context(), &pb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		response.GRPCError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}
