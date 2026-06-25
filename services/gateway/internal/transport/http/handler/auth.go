package handler

import (
	"net/http"

	"gateway/internal/config"
	"gateway/internal/transport/http/cookie"
	"gateway/internal/transport/http/dto"
	"gateway/internal/transport/http/helpers"
	"gateway/internal/transport/http/middleware"
	"gateway/internal/transport/http/response"
	"gateway/internal/validation"

	pb "github.com/hardsmile98/messager/sdk/auth/v1"
)

type Auth struct {
	client       pb.AuthServiceClient
	cookieConfig cookie.Settings
}

func NewAuth(client pb.AuthServiceClient, cfg *config.Config) *Auth {
	return &Auth{
		client: client,
		cookieConfig: cookie.Settings{
			Secure: cfg.CookieSecure,
			Domain: cfg.CookieDomain,
		},
	}
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

	cookie.SetAuthTokens(w, resp.AccessToken, resp.RefreshToken, cookie.TokenExpiry{
		AccessExpiresAt:  resp.AccessTokenExpiresAt.AsTime(),
		RefreshExpiresAt: resp.RefreshTokenExpiresAt.AsTime(),
	}, h.cookieConfig)

	response.JSON(w, http.StatusOK, map[string]string{
		"user_id": resp.UserId,
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

	cookie.SetAuthTokens(w, resp.AccessToken, resp.RefreshToken, cookie.TokenExpiry{
		AccessExpiresAt:  resp.AccessTokenExpiresAt.AsTime(),
		RefreshExpiresAt: resp.RefreshTokenExpiresAt.AsTime(),
	}, h.cookieConfig)

	response.JSON(w, http.StatusOK, map[string]string{
		"user_id": resp.UserId,
	})
}

func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, ok := middleware.RefreshTokenFromContext(r)

	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]string{
			"error": "session required",
		})
		return
	}

	resp, err := h.client.Logout(r.Context(), &pb.LogoutRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		response.GRPCError(w, err)
		return
	}

	cookie.ClearAuthTokens(w, h.cookieConfig)

	response.JSON(w, http.StatusOK, map[string]bool{
		"success": resp.Success,
	})
}

func (h *Auth) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, ok := middleware.RefreshTokenFromContext(r)

	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]string{
			"error": "session required",
		})
		return
	}

	resp, err := h.client.RefreshToken(r.Context(), &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		response.GRPCError(w, err)
		return
	}

	cookie.SetAuthTokens(w, resp.AccessToken, resp.RefreshToken, cookie.TokenExpiry{
		AccessExpiresAt:  resp.AccessTokenExpiresAt.AsTime(),
		RefreshExpiresAt: resp.RefreshTokenExpiresAt.AsTime(),
	}, h.cookieConfig)

	response.JSON(w, http.StatusOK, map[string]bool{
		"success": true,
	})
}
