//contains the handlers for login and callback
//this also connects OAuth success to application-session creation and adds logout.

package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ABHIJNA18/strava-ai-coach/internal/auth"
	"github.com/ABHIJNA18/strava-ai-coach/internal/strava"
)

// define struct
type AuthHandlers struct {
	DB           *sql.DB
	ClientID     string
	ClientSecret string
	Sessions     *auth.SessionManager
	SyncService  *strava.SyncService
	RedirectURI  string
	Secure       bool
}

// func to return handler
func NewAuthHandler(
	db *sql.DB,
	clientID string,
	clientSecret string,
	redirectURI string,
	sessions *auth.SessionManager,
	syncService *strava.SyncService,
	secure bool,
) *AuthHandlers {
	return &AuthHandlers{
		DB:           db,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Sessions:     sessions,
		SyncService:  syncService,
		Secure:       secure,
	}
}

// methods
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {

	// strava.LoginHandler(h.clientID) returns func and (w,r) is used to execute / call the function
	strava.LoginHandler(
		h.ClientID,
		h.RedirectURI,
		h.Secure,
	)(w, r)
}

func (h *AuthHandlers) Callback(w http.ResponseWriter, r *http.Request) {

	// strava.CallBackhandler() returns func and (w,r) is used to execute / call the function
	strava.CallbackHandler(
		h.ClientID,
		h.ClientSecret,
		h.RedirectURI,
		h.DB,
		h.Sessions,
		h.SyncService,
		h.Secure,
	)(w, r)
}

// Logout revokes only the current browser session, clears its cookie,
// and redirects the user to the public landing page.

func (h *AuthHandlers) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	cookie, err := r.Cookie(
		h.Sessions.CookieName(),
	)

	if err == nil {
		// Revoke only the session represented by this cookie.
		// Errors are intentionally ignored so logout remains graceful
		// for invalid or already-revoked cookies.
		_ = h.Sessions.RevokeSession(cookie.Value)
	}

	// The server clears the HttpOnly cookie.
	h.Sessions.ClearCookie(w)

	http.Redirect(
		w,
		r,
		"/",
		http.StatusFound,
	)
}
