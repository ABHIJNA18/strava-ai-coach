package strava

type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
}

type Athlete struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
}

type TokenResponse struct {
	TokenType    string  `json:"token_type"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresAt    int64   `json:"expires_at"`
	Athlete      Athlete `json:"athlete"`
}

type Activity struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Distance   float64 `json:"distance"`
	MovingTime int     `json:"moving_time"`
	Type       string  `json:"type"`
	StartDate  string  `json:"start_date"`
}
