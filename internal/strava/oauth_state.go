// This file protects the OAuth flow with a short-lived, single-use state cookie.
// generates, stores, validates, and clears the temporary OAuth CSRF state.

package strava

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"time"
)

const (
	oauthStateCookieName = "oauth_state"
	oauthStateLifetime   = 5 * time.Minute
)

func GenerateOAuthState() (string, error) {
	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func SetOAuthStateCookie(
	w http.ResponseWriter,
	state string,
	secure bool,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     oauthStateCookieName,
			Value:    state,
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(oauthStateLifetime.Seconds()),
			Expires:  time.Now().Add(oauthStateLifetime),
		},
	)
}

func ClearOAuthStateCookie(
	w http.ResponseWriter,
	secure bool,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     oauthStateCookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
			Expires:  time.Unix(1, 0),
		},
	)
}

func ValidateOAuthState(
	r *http.Request,
	callbackState string,
) bool {
	cookie, err := r.Cookie(oauthStateCookieName)
	if err != nil {
		return false
	}

	if callbackState == "" || cookie.Value == "" {
		return false
	}

	return subtle.ConstantTimeCompare(
		[]byte(cookie.Value),
		[]byte(callbackState),
	) == 1
}