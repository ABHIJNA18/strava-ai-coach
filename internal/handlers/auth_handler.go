package handlers

import (
	"database/sql"
	"net/http"
	"github.com/ABHIJNA18/strava-ai-coach/internal/strava"
)

//define struct
type AuthHandlers struct {
	DB 				*sql.DB
	ClientID 		string
	ClientSecret 	string
}

//func to return handler
func NewAuthHandler(db *sql.DB, clientID string,clientSecret string,) *AuthHandlers{
	return &AuthHandlers{
		DB : db,
		ClientID: clientID,
		ClientSecret: clientSecret,
	}
}

//methods 
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request){

	// strava.LoginHandler(h.clientID) returns func and (w,r) is used to execute / call the function
	strava.LoginHandler(
		h.ClientID,
	)(w, r)
}

func (h *AuthHandlers) Callback (w http.ResponseWriter, r *http.Request){

	// strava.CallBackhandler() returns func and (w,r) is used to execute / call the function
	strava.CallbackHandler(
		h.ClientID,
		h.ClientSecret,
		h.DB,
		)(w,r)
} 
