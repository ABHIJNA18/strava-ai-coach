// This file tests authentication middleware and athlete identity isolation.

package middleware

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/auth"
	"github.com/DATA-DOG/go-sqlmock"
)

func hashForTest(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func newTestSessionManager(
	t *testing.T,
) (
	*auth.SessionManager,
	sqlmock.Sqlmock,
) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"failed to create sqlmock database: %v",
			err,
		)
	}

	t.Cleanup(func() {
		db.Close()
	})

	sessionManager := auth.NewSessionManager(
		db,
		"session_id",
		false,
	)

	return sessionManager, mock
}

func expectValidSession(
	mock sqlmock.Sqlmock,
	rawToken string,
	sessionID int64,
	athleteID int64,
) {
	selectPattern := `(?s)SELECT.*FROM sessions.*WHERE token_hash = \$1`

	// Keep the fixture inside the session idle and absolute lifetimes.
	createdAt := time.Now().Add(-time.Hour)

	lastSeenAt := createdAt

	mock.ExpectQuery(selectPattern).
		WithArgs(hashForTest(rawToken)).
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
				sessionID,
				athleteID,
				hashForTest(rawToken),
				createdAt,
				lastSeenAt,
				nil,
			),
		)

	touchPattern := `(?s)UPDATE sessions.*SET last_seen_at`

	mock.ExpectExec(touchPattern).
		WithArgs(sessionID).
		WillReturnResult(
			sqlmock.NewResult(
				sessionID,
				1,
			),
		)
}

func expectRevokedSession(
	mock sqlmock.Sqlmock,
	rawToken string,
	sessionID int64,
	athleteID int64,
) {
	selectPattern := `(?s)SELECT.*FROM sessions.*WHERE token_hash = \$1`

	revokedAt := time.Date(
		2026,
		8,
		31,
		11,
		0,
		0,
		0,
		time.UTC,
	)

	mock.ExpectQuery(selectPattern).
		WithArgs(hashForTest(rawToken)).
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
				sessionID,
				athleteID,
				hashForTest(rawToken),
				revokedAt,
				revokedAt,
				revokedAt,
			),
		)
}

func TestRequireAuthValidSession(
	t *testing.T,
) {
	sessionManager, mock := newTestSessionManager(t)

	rawToken := "user-a-token"
	expectedAthleteID := int64(42)

	expectValidSession(
		mock,
		rawToken,
		1,
		expectedAthleteID,
	)

	handlerCalled := false

	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true

			athleteID, ok := AthleteIDFromContext(
				r.Context(),
			)

			if !ok {
				t.Fatal(
					"athlete ID was not present in context",
				)
			}

			if athleteID != expectedAthleteID {
				t.Fatalf(
					"expected athlete ID %d, got %d",
					expectedAthleteID,
					athleteID,
				)
			}

			w.WriteHeader(http.StatusOK)
		},
	)

	protected := RequireAuth(
		sessionManager,
		next,
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/dashboard",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session_id",
			Value: rawToken,
		},
	)

	recorder := httptest.NewRecorder()

	protected.ServeHTTP(
		recorder,
		request,
	)

	if !handlerCalled {
		t.Fatal(
			"protected handler was not called",
		)
	}

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unmet SQL expectation: %v",
			err,
		)
	}
}

func TestRequireAuthMissingCookie(
	t *testing.T,
) {
	sessionManager, mock := newTestSessionManager(t)

	handlerCalled := false

	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		},
	)

	protected := RequireAuth(
		sessionManager,
		next,
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/dashboard",
		nil,
	)

	recorder := httptest.NewRecorder()

	protected.ServeHTTP(
		recorder,
		request,
	)

	if handlerCalled {
		t.Fatal(
			"protected handler was called without a cookie",
		)
	}

	if recorder.Code != http.StatusFound {
		t.Fatalf(
			"expected redirect status 302, got %d",
			recorder.Code,
		)
	}

	if location := recorder.Header().Get("Location"); location != "/" {
		t.Fatalf(
			"expected redirect to /, got %s",
			location,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unexpected SQL call: %v",
			err,
		)
	}
}

func TestRequireAuthInvalidToken(
	t *testing.T,
) {
	sessionManager, mock := newTestSessionManager(t)

	rawToken := "nonexistent-token"

	selectPattern := `(?s)SELECT.*FROM sessions.*WHERE token_hash = \$1`

	mock.ExpectQuery(selectPattern).
		WithArgs(hashForTest(rawToken)).
		WillReturnError(sql.ErrNoRows)

	handlerCalled := false

	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		},
	)

	protected := RequireAuth(
		sessionManager,
		next,
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/dashboard",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session_id",
			Value: rawToken,
		},
	)

	recorder := httptest.NewRecorder()

	protected.ServeHTTP(
		recorder,
		request,
	)

	if handlerCalled {
		t.Fatal(
			"protected handler was called with an invalid token",
		)
	}

	if recorder.Code != http.StatusFound {
		t.Fatalf(
			"expected redirect status 302, got %d",
			recorder.Code,
		)
	}

	if location := recorder.Header().Get("Location"); location != "/" {
		t.Fatalf(
			"expected redirect to /, got %s",
			location,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unmet SQL expectation: %v",
			err,
		)
	}
}

func TestRequireAuthRevokedSession(
	t *testing.T,
) {
	sessionManager, mock := newTestSessionManager(t)

	rawToken := "revoked-token"

	expectRevokedSession(
		mock,
		rawToken,
		2,
		99,
	)

	handlerCalled := false

	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		},
	)

	protected := RequireAuth(
		sessionManager,
		next,
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/dashboard",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session_id",
			Value: rawToken,
		},
	)

	recorder := httptest.NewRecorder()

	protected.ServeHTTP(
		recorder,
		request,
	)

	if handlerCalled {
		t.Fatal(
			"protected handler was called with a revoked session",
		)
	}

	if recorder.Code != http.StatusFound {
		t.Fatalf(
			"expected redirect status 302, got %d",
			recorder.Code,
		)
	}

	if location := recorder.Header().Get("Location"); location != "/" {
		t.Fatalf(
			"expected redirect to /, got %s",
			location,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unmet SQL expectation: %v",
			err,
		)
	}
}

func TestUserSessionsResolveToDifferentAthletes(
	t *testing.T,
) {
	tests := []struct {
		name      string
		token     string
		athleteID int64
		sessionID int64
	}{
		{
			name:      "user A",
			token:     "user-a-token",
			athleteID: 42,
			sessionID: 1,
		},
		{
			name:      "user B",
			token:     "user-b-token",
			athleteID: 84,
			sessionID: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sessionManager, mock := newTestSessionManager(t)

			expectValidSession(
				mock,
				test.token,
				test.sessionID,
				test.athleteID,
			)

			var receivedAthleteID int64

			next := http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					var ok bool

					receivedAthleteID, ok =
						AthleteIDFromContext(
							r.Context(),
						)

					if !ok {
						t.Fatal(
							"athlete ID missing from context",
						)
					}

					w.WriteHeader(http.StatusOK)
				},
			)

			protected := RequireAuth(
				sessionManager,
				next,
			)

			request := httptest.NewRequest(
				http.MethodGet,
				"/activities?athlete_id=999999",
				nil,
			)

			request.AddCookie(
				&http.Cookie{
					Name:  "session_id",
					Value: test.token,
				},
			)

			recorder := httptest.NewRecorder()

			protected.ServeHTTP(
				recorder,
				request,
			)

			if receivedAthleteID != test.athleteID {
				t.Fatalf(
					"expected athlete ID %d, got %d",
					test.athleteID,
					receivedAthleteID,
				)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf(
					"unmet SQL expectation: %v",
					err,
				)
			}
		})
	}
}
