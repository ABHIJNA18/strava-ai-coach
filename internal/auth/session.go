// This file creates and validates persistent application authentication cookies and sends it to the client browser

package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
)

var ErrInvalidSession = errors.New(
	"invalid application session",
)

const (
	sessionAbsoluteLifetime = 30 * 24 * time.Hour
	sessionIdleTimeout      = 7 * 24 * time.Hour
)

type SessionManager struct {
	db         *sql.DB
	cookieName string
	secure     bool
}

func NewSessionManager(
	db *sql.DB,
	cookieName string,
	secure bool,
) *SessionManager {
	return &SessionManager{
		db:         db,
		cookieName: cookieName,
		secure:     secure,
	}
}

func generateToken() (string, error) {
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}

func (m *SessionManager) CreateSession(
	athleteID int64,
) (string, error) {
	rawToken, err := generateToken()
	if err != nil {
		return "", err
	}

	session := database.Session{
		AthleteID: athleteID,
		TokenHash: hashToken(rawToken),
	}

	_, err = database.CreateSession(
		m.db,
		session,
	)
	if err != nil {
		return "", err
	}

	//return raw token and hashed token to be stored in database.
	return rawToken, nil
}

func (m *SessionManager) LoadSession(
	rawToken string,
) (*database.Session, error) {

	if rawToken == "" {
		return nil, ErrInvalidSession
	}

	session, err := database.GetSessionByTokenHash(
		m.db,
		hashToken(rawToken),
	)
	if err != nil {
		return nil, ErrInvalidSession
	}

	if session.RevokedAt != nil {
		return nil, ErrInvalidSession
	}

	//check session expiration based on absolute lifetime and idle timeout
	now := time.Now()

	if now.Sub(session.CreatedAt) > sessionAbsoluteLifetime {
		return nil, ErrInvalidSession
	}

	if now.Sub(session.LastSeenAt) > sessionIdleTimeout {
		return nil, ErrInvalidSession
	}

	//update the last seen time of the session in database
	err = database.TouchSession(
		m.db,
		session.ID,
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (m *SessionManager) RevokeSession(
	rawToken string,
) error {
	if rawToken == "" {
		return nil
	}

	session, err := database.GetSessionByTokenHash(
		m.db,
		hashToken(rawToken),
	)
	if err != nil {
		return nil
	}

	return database.RevokeSession(
		m.db,
		session.ID,
	)
}

// Setting the cookie which is sent to the client browser, which is the raw token
func (m *SessionManager) SetCookie(
	w http.ResponseWriter,
	rawToken string,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     m.cookieName,
			Value:    rawToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   m.secure,
			SameSite: http.SameSiteLaxMode,
		},
	)
}

func (m *SessionManager) ClearCookie(
	w http.ResponseWriter,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     m.cookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   m.secure,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		},
	)
}

func (m *SessionManager) CookieName() string {
	return m.cookieName
}
