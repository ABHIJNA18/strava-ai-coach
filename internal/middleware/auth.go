
// This file protects HTTP routes using the application authentication cookie.

package middleware

import (
	"context"
	"net/http"

	"github.com/ABHIJNA18/strava-ai-coach/internal/auth"
)

type contextKey string

const athleteIDContextKey contextKey = "athleteID"

func RequireAuth(
	sessionManager *auth.SessionManager,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(
				sessionManager.CookieName(),
			)

			//if cookie doesn't exist, redirect to login page
			if err != nil {
				http.Redirect(
					w,
					r,
					"/",
					http.StatusFound,
				)
				return
			}

			//if cookie exists, get the session details using cookie value which is the raw token

			session, err := sessionManager.LoadSession(
				cookie.Value,
			)
			if err != nil {
				sessionManager.ClearCookie(w)

				http.Redirect(
					w,
					r,
					"/",
					http.StatusFound,
				)
				return
			}
			
			//add additional context to thr request, which is the athleteID from the session,this information goes to the next handler 
			//“Create a new context and store session.AthleteID under the key athleteIDContextKey
			//You are doing the connection explicitly at context.WithValue.
			ctx := context.WithValue(
				r.Context(),
				athleteIDContextKey,
				session.AthleteID,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		},
	)
}

// retrives the athlete ID from the request context, if it exists.

func AthleteIDFromContext(
	ctx context.Context,
) (int64, bool) {
	athleteID, ok := ctx.Value(
		athleteIDContextKey,
	).(int64)

	return athleteID, ok
}