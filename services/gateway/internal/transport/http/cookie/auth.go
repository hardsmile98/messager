package cookie

import (
	"net/http"
	"time"
)

const (
	AccessTokenName  = "access_token"
	RefreshTokenName = "refresh_token"

	refreshTokenPath = "/api/v1/auth"
)

type Settings struct {
	Secure bool
	Domain string
}

type TokenExpiry struct {
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

func SetAuthTokens(w http.ResponseWriter, accessToken, refreshToken string, expiry TokenExpiry, s Settings) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenName,
		Value:    accessToken,
		Path:     "/",
		MaxAge:   maxAgeSeconds(expiry.AccessExpiresAt),
		HttpOnly: true,
		Secure:   s.Secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   s.Domain,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenName,
		Value:    refreshToken,
		Path:     refreshTokenPath,
		MaxAge:   maxAgeSeconds(expiry.RefreshExpiresAt),
		HttpOnly: true,
		Secure:   s.Secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   s.Domain,
	})
}

func maxAgeSeconds(expiresAt time.Time) int {
	seconds := int(time.Until(expiresAt).Seconds())
	if seconds < 0 {
		return 0
	}

	return seconds
}

func ClearAuthTokens(w http.ResponseWriter, s Settings) {
	clear := func(name, path string) {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   s.Secure,
			SameSite: http.SameSiteLaxMode,
			Domain:   s.Domain,
		})
	}

	clear(AccessTokenName, "/")
	clear(RefreshTokenName, refreshTokenPath)
}

func AccessToken(r *http.Request) string {
	if c, err := r.Cookie(AccessTokenName); err == nil {
		return c.Value
	}

	return ""
}

func RefreshToken(r *http.Request) string {
	if c, err := r.Cookie(RefreshTokenName); err == nil {
		return c.Value
	}

	return ""
}
