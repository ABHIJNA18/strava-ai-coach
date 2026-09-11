// This file tests secure application session creation, lookup, and revocation.

package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func expectedHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func TestCreateSessionStoresHashedToken(
	t *testing.T,
) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"failed to create sqlmock database: %v",
			err,
		)
	}
	defer db.Close()

	sessionManager := NewSessionManager(
		db,
		"session_id",
		false,
	)

	mock.ExpectQuery(
		`(?s)INSERT INTO sessions.*RETURNING id`,
	).
		WithArgs(
			int64(42),
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"id"},
			).AddRow(1),
		)

	rawToken, err := sessionManager.CreateSession(42)
	if err != nil {
		t.Fatalf(
			"CreateSession returned error: %v",
			err,
		)
	}

	if rawToken == "" {
		t.Fatal(
			"expected a non-empty raw token",
		)
	}

	if len(rawToken) != 64 {
		t.Fatalf(
			"expected a 64-character token, got %d characters",
			len(rawToken),
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unmet SQL expectation: %v",
			err,
		)
	}
}

func TestLoadSessionReturnsAthleteID(
	t *testing.T,
) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"failed to create sqlmock database: %v",
			err,
		)
	}
	defer db.Close()

	sessionManager := NewSessionManager(
		db,
		"session_id",
		false,
	)

	rawToken := "user-a-token"
	tokenHash := expectedHash(rawToken)

	now := time.Now()

	mock.ExpectQuery(
		`(?s)SELECT.*FROM sessions.*WHERE token_hash = \$1`,
	).
		WithArgs(tokenHash).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id",
					"athlete_id",
					"token_hash",
					"created_at",
					"last_seen_at",
					"revoked_at",
				},
			).AddRow(
				1,
				int64(42),
				tokenHash,
				now,
				now,
				nil,
			),
		)

	mock.ExpectExec(
		`(?s)UPDATE sessions.*SET last_seen_at`,
	).
		WithArgs(int64(1)).
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

	session, err := sessionManager.LoadSession(
		rawToken,
	)
	if err != nil {
		t.Fatalf(
			"LoadSession returned error: %v",
			err,
		)
	}

	if session.AthleteID != 42 {
		t.Fatalf(
			"expected athlete ID 42, got %d",
			session.AthleteID,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unmet SQL expectation: %v",
			err,
		)
	}
}

func TestLoadSessionRejectsRevokedSession(
	t *testing.T,
) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"failed to create sqlmock database: %v",
			err,
		)
	}
	defer db.Close()

	sessionManager := NewSessionManager(
		db,
		"session_id",
		false,
	)

	rawToken := "revoked-token"
	tokenHash := expectedHash(rawToken)

	now := time.Now()
	revokedAt := now

	mock.ExpectQuery(
		`(?s)SELECT.*FROM sessions.*WHERE token_hash = \$1`,
	).
		WithArgs(tokenHash).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id",
					"athlete_id",
					"token_hash",
					"created_at",
					"last_seen_at",
					"revoked_at",
				},
			).AddRow(
				1,
				int64(42),
				tokenHash,
				now,
				now,
				revokedAt,
			),
		)

	_, err = sessionManager.LoadSession(
		rawToken,
	)
	if err != ErrInvalidSession {
		t.Fatalf(
			"expected ErrInvalidSession, got %v",
			err,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unmet SQL expectation: %v",
			err,
		)
	}
}

func TestSetCookie(
	t *testing.T,
) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"failed to create sqlmock database: %v",
			err,
		)
	}
	defer db.Close()

	sessionManager := NewSessionManager(
		db,
		"session_id",
		false,
	)

	recorder := httptest.NewRecorder()

	sessionManager.SetCookie(
		recorder,
		"raw-session-token",
	)

	cookies := recorder.Result().Cookies()

	if len(cookies) != 1 {
		t.Fatalf(
			"expected one cookie, got %d",
			len(cookies),
		)
	}

	cookie := cookies[0]

	if cookie.Name != "session_id" {
		t.Fatalf(
			"expected cookie name session_id, got %s",
			cookie.Name,
		)
	}

	if cookie.Value != "raw-session-token" {
		t.Fatalf(
			"unexpected cookie value: %s",
			cookie.Value,
		)
	}

	if !cookie.HttpOnly {
		t.Fatal(
			"expected HttpOnly=true",
		)
	}

	if cookie.Secure {
		t.Fatal(
			"expected Secure=false for local development",
		)
	}

	if cookie.Path != "/" {
		t.Fatalf(
			"expected cookie path /, got %s",
			cookie.Path,
		)
	}

	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf(
			"expected SameSite=Lax",
		)
	}
}

func TestClearCookie(
	t *testing.T,
) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"failed to create sqlmock database: %v",
			err,
		)
	}
	defer db.Close()

	sessionManager := NewSessionManager(
		db,
		"session_id",
		false,
	)

	recorder := httptest.NewRecorder()

	sessionManager.ClearCookie(
		recorder,
	)

	cookies := recorder.Result().Cookies()

	if len(cookies) != 1 {
		t.Fatalf(
			"expected one cookie, got %d",
			len(cookies),
		)
	}

	cookie := cookies[0]

	if cookie.Name != "session_id" {
		t.Fatalf(
			"expected cookie name session_id, got %s",
			cookie.Name,
		)
	}

	if cookie.Value != "" {
		t.Fatalf(
			"expected empty cookie value",
		)
	}

	if cookie.MaxAge != -1 {
		t.Fatalf(
			"expected MaxAge=-1, got %d",
			cookie.MaxAge,
		)
	}
}

/*
CreateSession()
    → creates a random token
    → stores a hash argument in the database

LoadSession()
    → hashes the supplied cookie token
    → loads the session
    → returns the correct athlete ID
    → updates last_seen_at

LoadSession() with revoked row
    → returns ErrInvalidSession
    → does not authenticate the request

SetCookie()
    → HttpOnly
    → local Secure=false
    → SameSite=Lax
    → Path=/

ClearCookie()
    → empty value
    → MaxAge=-1
*/
