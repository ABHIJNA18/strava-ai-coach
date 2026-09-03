// This file tests logout session revocation and cookie clearing.
//this verifies that logout revokes the database session and clears the browser cookie.

package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"database/sql"

	"github.com/ABHIJNA18/strava-ai-coach/internal/auth"
	"github.com/DATA-DOG/go-sqlmock"
)

func testTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func TestLogoutRevokesSessionAndClearsCookie(
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

	sessionManager := auth.NewSessionManager(
		db,
		"session_id",
		false,
	)

	authHandler := &AuthHandlers{
		Sessions: sessionManager,
	}

	rawToken := "logout-token"

	createdAt := time.Date(
		2026,
		8,
		31,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	selectPattern := `(?s)SELECT.*FROM sessions.*WHERE token_hash = \$1`

	mock.ExpectQuery(selectPattern).
		WithArgs(testTokenHash(rawToken)).
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
				10,
				42,
				testTokenHash(rawToken),
				createdAt,
				createdAt,
				nil,
			),
		)

	revokePattern := `(?s)UPDATE sessions.*SET revoked_at`

	mock.ExpectExec(revokePattern).
		WithArgs(int64(10)).
		WillReturnResult(
			sqlmock.NewResult(
				10,
				1,
			),
		)

	request := httptest.NewRequest(
		http.MethodPost,
		"/logout",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session_id",
			Value: rawToken,
		},
	)

	recorder := httptest.NewRecorder()

	authHandler.Logout(
		recorder,
		request,
	)

	if recorder.Code != http.StatusFound {
		t.Fatalf(
			"expected status 302, got %d",
			recorder.Code,
		)
	}

	if location := recorder.Header().Get("Location"); location != "/" {
		t.Fatalf(
			"expected redirect to /, got %s",
			location,
		)
	}

	cookies := recorder.Result().Cookies()

	if len(cookies) == 0 {
		t.Fatal(
			"expected a clearing Set-Cookie header",
		)
	}

	clearedCookie := cookies[0]

	if clearedCookie.Name != "session_id" {
		t.Fatalf(
			"expected cookie session_id, got %s",
			clearedCookie.Name,
		)
	}

	if clearedCookie.MaxAge != -1 {
		t.Fatalf(
			"expected MaxAge -1, got %d",
			clearedCookie.MaxAge,
		)
	}

	if clearedCookie.Value != "" {
		t.Fatalf(
			"expected empty cookie value, got %q",
			clearedCookie.Value,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unmet SQL expectation: %v",
			err,
		)
	}
}

func TestLogoutWithoutCookieStillClearsCookie(
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

	sessionManager := auth.NewSessionManager(
		db,
		"session_id",
		false,
	)

	authHandler := &AuthHandlers{
		Sessions: sessionManager,
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/logout",
		nil,
	)

	recorder := httptest.NewRecorder()

	authHandler.Logout(
		recorder,
		request,
	)

	if recorder.Code != http.StatusFound {
		t.Fatalf(
			"expected status 302, got %d",
			recorder.Code,
		)
	}

	if len(recorder.Result().Cookies()) == 0 {
		t.Fatal(
			"expected cookie clearing response",
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unexpected SQL call: %v",
			err,
		)
	}
}

func TestLogoutGETIsRejectedWithoutRevoking(
	t *testing.T,
) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sessionManager := auth.NewSessionManager(
		db,
		"session_id",
		false,
	)

	authHandler := &AuthHandlers{
		Sessions: sessionManager,
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/logout",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session_id",
			Value: "valid-token",
		},
	)

	recorder := httptest.NewRecorder()

	authHandler.Logout(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected 405, got %d",
			recorder.Code,
		)
	}

	if len(recorder.Result().Cookies()) != 0 {
		t.Fatal("GET /logout must not clear or modify the cookie")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"GET /logout unexpectedly accessed the database: %v",
			err,
		)
	}
}

func TestLogoutInvalidSessionStillClearsCookie(
	t *testing.T,
) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sessionManager := auth.NewSessionManager(
		db,
		"session_id",
		false,
	)

	authHandler := &AuthHandlers{
		Sessions: sessionManager,
	}

	rawToken := "invalid-token"

	selectPattern := `(?s)SELECT.*FROM sessions.*WHERE token_hash = \$1`

	mock.ExpectQuery(selectPattern).
		WithArgs(testTokenHash(rawToken)).
		WillReturnError(sql.ErrNoRows)

	request := httptest.NewRequest(
		http.MethodPost,
		"/logout",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session_id",
			Value: rawToken,
		},
	)

	recorder := httptest.NewRecorder()

	authHandler.Logout(recorder, request)

	if recorder.Code != http.StatusFound {
		t.Fatalf(
			"expected redirect status 302, got %d",
			recorder.Code,
		)
	}

	if recorder.Header().Get("Location") != "/" {
		t.Fatalf("expected redirect to /")
	}

	cookies := recorder.Result().Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected a clearing cookie")
	}

	if cookies[0].MaxAge != -1 {
		t.Fatalf(
			"expected MaxAge -1, got %d",
			cookies[0].MaxAge,
		)
	}

	if cookies[0].Value != "" {
		t.Fatalf(
			"expected empty cookie value, got %q",
			cookies[0].Value,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unexpected SQL behavior: %v",
			err,
		)
	}
}

func TestLogoutRevokesOnlyCurrentSession(
	t *testing.T,
) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sessionManager := auth.NewSessionManager(
		db,
		"session_id",
		false,
	)

	authHandler := &AuthHandlers{
		Sessions: sessionManager,
	}

	rawToken := "user-a-current-session"
	tokenHash := testTokenHash(rawToken)

	createdAt := time.Date(
		2026,
		8,
		31,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	selectPattern := `(?s)SELECT.*FROM sessions.*WHERE token_hash = \$1`

	mock.ExpectQuery(selectPattern).
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
				int64(10),
				int64(42),
				tokenHash,
				createdAt,
				createdAt,
				nil,
			),
		)

	revokePattern := `(?s)UPDATE sessions.*SET revoked_at`

	mock.ExpectExec(revokePattern).
		WithArgs(int64(10)).
		WillReturnResult(
			sqlmock.NewResult(10, 1),
		)

	request := httptest.NewRequest(
		http.MethodPost,
		"/logout",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session_id",
			Value: rawToken,
		},
	)

	recorder := httptest.NewRecorder()

	authHandler.Logout(recorder, request)

	if recorder.Code != http.StatusFound {
		t.Fatalf(
			"expected status 302, got %d",
			recorder.Code,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unexpected SQL behavior: %v",
			err,
		)
	}
}