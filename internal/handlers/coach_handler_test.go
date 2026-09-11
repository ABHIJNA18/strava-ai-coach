// This file tests that coaching handlers use session identity instead of client input.
//this tests that the handler uses the authenticated athlete ID from context and ignores a client-supplied athlete ID.

package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/auth"
	"github.com/ABHIJNA18/strava-ai-coach/internal/middleware"
	"github.com/DATA-DOG/go-sqlmock"
)

type fakeCoachService struct {
	receivedAthleteID int64
}

func (f *fakeCoachService) AnalyzeRecentRuns(
	ctx context.Context,
	athleteID int64,
) (string, error) {
	f.receivedAthleteID = athleteID
	return "test summary", nil
}

func (f *fakeCoachService) GenerateCoaching(
	ctx context.Context,
	athleteID int64,
	goal string,
) (string, error) {
	f.receivedAthleteID = athleteID
	return `{
		"observations": [],
		"progress": [],
		"risks": [],
		"recommendations": []
	}`, nil
}

func coachTestHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func expectCoachSession(
	mock sqlmock.Sqlmock,
	token string,
	sessionID int64,
	athleteID int64,
) {
	// Keep the fixture inside the session idle and absolute lifetimes.
	now := time.Now().Add(-time.Hour)

	selectPattern := `(?s)SELECT.*FROM sessions.*WHERE token_hash = \$1`

	mock.ExpectQuery(selectPattern).
		WithArgs(coachTestHash(token)).
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
				coachTestHash(token),
				now,
				now,
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

func TestGetReportUsesAuthenticatedAthlete(
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

	fakeService := &fakeCoachService{}

	handler := NewCoachHandler(
		fakeService,
	)

	rawToken := "user-a-token"
	expectedAthleteID := int64(42)

	expectCoachSession(
		mock,
		rawToken,
		1,
		expectedAthleteID,
	)

	protectedHandler := middleware.RequireAuth(
		sessionManager,
		http.HandlerFunc(
			handler.GetReport,
		),
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/coach/report?athlete_id=999",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session_id",
			Value: rawToken,
		},
	)

	recorder := httptest.NewRecorder()

	protectedHandler.ServeHTTP(
		recorder,
		request,
	)

	if fakeService.receivedAthleteID != expectedAthleteID {
		t.Fatalf(
			"expected athlete ID %d, got %d",
			expectedAthleteID,
			fakeService.receivedAthleteID,
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

func TestDifferentUsersRemainIsolated(
	t *testing.T,
) {
	tests := []struct {
		name      string
		token     string
		sessionID int64
		athleteID int64
	}{
		{
			name:      "user A",
			token:     "user-a-token",
			sessionID: 1,
			athleteID: 42,
		},
		{
			name:      "user B",
			token:     "user-b-token",
			sessionID: 2,
			athleteID: 84,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

			fakeService := &fakeCoachService{}

			handler := NewCoachHandler(
				fakeService,
			)

			expectCoachSession(
				mock,
				test.token,
				test.sessionID,
				test.athleteID,
			)

			protectedHandler := middleware.RequireAuth(
				sessionManager,
				http.HandlerFunc(
					handler.GetReport,
				),
			)

			request := httptest.NewRequest(
				http.MethodGet,
				"/coach/report?athlete_id=999999",
				nil,
			)

			request.AddCookie(
				&http.Cookie{
					Name:  "session_id",
					Value: test.token,
				},
			)

			recorder := httptest.NewRecorder()

			protectedHandler.ServeHTTP(
				recorder,
				request,
			)

			if fakeService.receivedAthleteID != test.athleteID {
				t.Fatalf(
					"expected athlete ID %d, got %d",
					test.athleteID,
					fakeService.receivedAthleteID,
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
