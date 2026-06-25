package handler

import (
	"encoding/json"
	"net/http"

	pb "github.com/hardsmile98/messager/sdk/auth/v1"
)

type AuthHandler struct {
	client pb.AuthServiceClient
}

func NewAuthHandler(client pb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Device   string `json:"device"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	resp, err := h.client.Register(r.Context(), &pb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Device:   req.Device,
	})

	if err != nil {
		writeGrpcError(w, err)
		return
	}

	writeJson(w, http.StatusOK, map[string]string{
		"user_id":       resp.UserId,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Device   string `json:"device"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	resp, err := h.client.Login(r.Context(), &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
		Device:   req.Device,
	})

	if err != nil {
		writeGrpcError(w, err)
		return
	}

	writeJson(w, http.StatusOK, map[string]string{
		"user_id":       resp.UserId,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	resp, err := h.client.Logout(r.Context(), &pb.LogoutRequest{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		writeGrpcError(w, err)
		return
	}

	writeJson(w, http.StatusOK, map[string]bool{
		"success": resp.Success,
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	resp, err := h.client.RefreshToken(r.Context(), &pb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		writeGrpcError(w, err)
		return
	}

	writeJson(w, http.StatusOK, map[string]string{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}
