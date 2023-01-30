package oauth

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	OauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8000/oauth/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	OauthStateString = ""
)

func InitializeAuth() {
	file, err := os.ReadFile("./credentials/oauth2-credentials.json")

	if err != nil {
		panic(err)
	}

	var creds interface{}
	err = json.Unmarshal(file, &creds)

	if err != nil {
		panic(err)
	}

	OauthConf.ClientID = creds.(map[string]interface{})["web"].(map[string]interface{})["client_id"].(string)
	OauthConf.ClientSecret = creds.(map[string]interface{})["web"].(map[string]interface{})["client_secret"].(string)
	OauthStateString = creds.(map[string]interface{})["state_string"].(string)
}
