//the callback must create an application authentication cookie after the athlete and Strava tokens are saved.
//OAuth should perform the 90-day sync only when initial_sync_completed_at is still NULL

package strava

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/auth"
	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
)

// contains code which handles the login and callback routes for the OAuth flow
func LoginHandler(
	clientID string,
	redirectURI string,
	secure bool,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := GenerateOAuthState()
		if err != nil {
			http.Error(
				w,
				"failed to initialize OAuth login",
				http.StatusInternalServerError,
			)
			return
		}

		SetOAuthStateCookie(
			w,
			state,
			secure,
		)

		query := url.Values{}
		query.Set("client_id", clientID)
		query.Set("response_type", "code")
		query.Set("redirect_uri", redirectURI)
		query.Set("approval_prompt", "auto")
		query.Set("scope", "read,activity:read_all")
		query.Set("state", state)

		authURL := "https://www.strava.com/oauth/authorize?" +
			query.Encode()

		http.Redirect(
			w,
			r,
			authURL,
			http.StatusFound,
		)
	}
}

func CallbackHandler(
	clientID string,
	clientSecret string,
	redirectURI string,
	db *sql.DB,
	sessions *auth.SessionManager,
	syncService *SyncService,
	secure bool,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//check callback state with the state stored in the cookie to prevent CSRF attacks
		callbackState := r.URL.Query().Get("state")

		if !ValidateOAuthState(r, callbackState) {
			ClearOAuthStateCookie(w, secure)

			http.Error(
				w,
				"invalid OAuth state",
				http.StatusBadRequest,
			)
			return
		}

		ClearOAuthStateCookie(w, secure)

		//Oauth error handling
		oauthError := r.URL.Query().Get("error")
		if oauthError != "" {
			http.Error(w, "Oauth authorisation failed : "+oauthError, http.StatusBadRequest)
			fmt.Printf("OAuth authorization failed: %s\n", oauthError)
			return
		}
		//get authorsation code
		auth_code := r.URL.Query().Get("code")

		if auth_code == "" {
			http.Error(w, "Authorization code not found in the request", http.StatusBadRequest)
			fmt.Println("Authorization code not found in the request")
			return
		}
		fmt.Println("Recevived auth code")

		//exchange auth code for access token
		tokenResponse, err := ExchangeTokenForCode(clientID, clientSecret, auth_code, redirectURI)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		//get athelete details using access token
		athlete, err := GetAthlete(tokenResponse.AccessToken)
		if err != nil {
			http.Error(w, "failed to fetch athlete", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Athlete data fetched successfully %s, %s \n", athlete.Firstname, athlete.Lastname)

		// change from strava.Athlete to database.Athlete to store athlete in db
		dbAthlete := database.Athlete{

			StravaAthleteID: athlete.ID,
			Firstname:       athlete.Firstname,
			Lastname:        athlete.Lastname,
		}

		//save the athlete to db and return athlete ID
		athleteID, err := database.SaveAthlete(db, dbAthlete)
		if err != nil {
			fmt.Println("Failed to save athlete to databse :", err)
			return
		}

		//convert oauth token details to db's oauth model
		dbToken := database.OAuthToken{

			AthleteID:    athleteID,
			AccessToken:  tokenResponse.AccessToken,
			RefreshToken: tokenResponse.RefreshToken,
			ExpiresAt:    tokenResponse.ExpiresAt,
		}

		//after getting athlete data, save the oauth token details with db's athleteID
		err = database.SaveOauthToken(db, dbToken)
		if err != nil {
			fmt.Println("Save tokens to databse failed")
			return
		}

		fmt.Println("Tokens saved to DB sucessfully")

		//OAuth should perform the 90-day sync only when initial_sync_completed_at is still NULL
		initialSyncCompletedAt, err :=
			database.GetInitialSyncCompletedAt(
				db,
				athleteID,
			)
		if err != nil {
			http.Error(
				w,
				"failed to check initial synchronization state",
				http.StatusInternalServerError,
			)
			return
		}

		if initialSyncCompletedAt == nil {
			fmt.Println("Starting initial 90-day activity synchronization")

			after := time.Now().AddDate(0, 0, -90)

			if err := syncService.SyncActivities(
				athleteID,
				after,
			); err != nil {
				fmt.Println(
					"Initial activity synchronization failed:",
					err,
				)

				http.Error(
					w,
					"Initial activity synchronization failed. Please try again.",
					http.StatusInternalServerError,
				)
				return
			}

			if err := database.MarkInitialSyncCompleted(
				db,
				athleteID,
			); err != nil {
				http.Error(
					w,
					"failed to mark initial synchronization complete",
					http.StatusInternalServerError,
				)
				return
			}

			fmt.Println("Initial 90-day synchronization completed")
		} else {
			fmt.Println(
				"Initial synchronization already completed; skipping",
			)
		}

		//the callback must create an application authentication cookie after the athlete and Strava tokens are saved.

		sessionToken, err := sessions.CreateSession(
			athleteID,
		)
		if err != nil {
			http.Error(
				w,
				"Failed to create application session",
				http.StatusInternalServerError,
			)
			return
		}

		sessions.SetCookie(
			w,
			sessionToken,
		)

		// redirecting the users to dashboard page from login page after successfull login
		http.Redirect(w, r, "/dashboard", http.StatusFound)

	}
}

/*func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	fmt.Fprintf(w, "Authorization Code: %s", code)
}

this veriosn of the code also works
because CallbackHandler doesn’t currently need any configuration.
We can directly write the handler itself instead of a function that returns one like loginHandler
*/
