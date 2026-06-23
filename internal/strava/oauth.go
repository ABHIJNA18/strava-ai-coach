package strava

import (
	"fmt"
	"net/http"
)

// contains code which handles the login and callback routes for the OAuth flow
func LoginHandler(clientID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authURL := fmt.Sprintf("https://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=http://localhost:8080/oauth/callback&approval_prompt=force&scope=read,activity:read_all", clientID)
		http.Redirect(w, r, authURL, http.StatusFound)
	}
}
func CallbackHandler(clientID string, clientSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Oauth error handling
		oauthError := r.URL.Query().Get("error")
		if oauthError != ""{
			http.Error(w, "Oauth authorisation failed : " +oauthError, http.StatusBadRequest)
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

		//get activities using access token
		activities, err := GetActivities(tokenResponse.AccessToken)
		if err != nil {
			http.Error(w, "Failed to fetch activities", http.StatusInternalServerError)
			return
		}
		fmt.Printf(" %d Activities fetched successfully :", len(activities))
		fmt.Fprintf(w, " %d Activities fetched successfully :", len(activities))

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
