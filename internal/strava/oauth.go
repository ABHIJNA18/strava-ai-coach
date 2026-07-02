package strava

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
)

// contains code which handles the login and callback routes for the OAuth flow
func LoginHandler(clientID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authURL := fmt.Sprintf("https://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=http://localhost:8080/oauth/callback&approval_prompt=force&scope=read,activity:read_all", clientID)
		http.Redirect(w, r, authURL, http.StatusFound)
	}
}
func CallbackHandler(clientID string, clientSecret string, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
		tokenResponse, err := ExchangeTokenForCode(clientID, clientSecret, auth_code)
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

		//get all the activities using access token
		activities, err := GetAllActivities(tokenResponse.AccessToken)
		if err != nil {
			http.Error(w, "Failed to fetch all activities", http.StatusInternalServerError)
			return
		}
		fmt.Printf(" %d Activities fetched successfully :", len(activities))
		fmt.Fprintf(w, " %d Activities fetched successfully :", len(activities))

		//creating the variable of type database.Activities to insert into DB
		var dbActivities []database.Activity

		//parse timestamps and append
		for _, activity := range activities {

			startDate, err := time.Parse(time.RFC3339, activity.StartDate)
			if err != nil {
				return
			}

			startDateLocal, err := time.Parse(time.RFC3339, activity.StartDateLocal)
			if err != nil {
				return
			}

			dbActivities = append(
				dbActivities,
				database.Activity{

					StravaActivityID:     activity.ID,
					AthleteID:            athleteID,
					Name:                 activity.Name,
					Type:                 activity.Type,
					SportType:            activity.SportType,
					Distance:             activity.Distance,
					MovingTime:           activity.MovingTime,
					ElapsedTime:          activity.ElapsedTime,
					TotalElevationGain:   activity.TotalElevationGain,
					AverageSpeed:         activity.AverageSpeed,
					MaxSpeed:             activity.MaxSpeed,
					AverageHeartrate:     activity.AverageHeartrate,
					MaxHeartrate:         activity.MaxHeartrate,
					AverageCadence:       activity.AverageCadence,
					AverageWatts:         activity.AverageWatts,
					MaxWatts:             activity.MaxWatts,
					WeightedAverageWatts: activity.WeightedAverageWatts,
					Kilojoules:           activity.Kilojoules,
					SufferScore:          activity.SufferScore,
					DeviceName:           activity.DeviceName,
					StartDate:            startDate,
					StartDateLocal:       startDateLocal,
				},
			)

		}

		err = database.SaveActivities(db, dbActivities)

		if err != nil {
			fmt.Println(
				"Failed to save activities:",
				err,
			)
			return
		}

		fmt.Printf(
			"\n %d activities saved successfully\n",
			len(dbActivities),
		)

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
