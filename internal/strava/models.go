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

type RefreshTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
}

type Activity struct {
	ID                     int64   `json:"id"`

	Name                   string  `json:"name"`
	Type                   string  `json:"type"`
	SportType              string  `json:"sport_type"`

	Distance               float64 `json:"distance"`

	MovingTime             int     `json:"moving_time"`
	ElapsedTime            int     `json:"elapsed_time"`

	TotalElevationGain     float64 `json:"total_elevation_gain"`

	AverageSpeed           float64 `json:"average_speed"`
	MaxSpeed               float64 `json:"max_speed"`

	AverageHeartrate       float64 `json:"average_heartrate"`
	MaxHeartrate           float64 `json:"max_heartrate"`

	AverageCadence         float64 `json:"average_cadence"`

	AverageWatts           float64 `json:"average_watts"`
	MaxWatts               float64 `json:"max_watts"`
	WeightedAverageWatts   float64 `json:"weighted_average_watts"`

	Kilojoules             float64 `json:"kilojoules"`

	SufferScore            float64 `json:"suffer_score"`

	DeviceName             string  `json:"device_name"`

	StartDate              string  `json:"start_date"`
	StartDateLocal         string  `json:"start_date_local"`
}
