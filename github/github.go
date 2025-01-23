package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// GitHub OAuth credentials
var (
	clientID     = ""
	clientSecret = ""
	oauthURL     = "https://github.com/login/oauth"
	apiURL       = "https://api.github.com/user"
	redirectURI  = "http://localhost:5000/github/callback"
)

func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/github/callback", handleCallback)

	fmt.Println("Server running on http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

// handleLogin redirects the user to GitHub's OAuth authorization URL.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/authorize?client_id=%s&redirect_uri=%s&scope=user", oauthURL, clientID, redirectURI)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleCallback handles the OAuth callback from GitHub.
func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
		return
	}

	token, err := getAccessToken(code)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		log.Println("Error fetching access token:", err)
		return
	}

	user, err := getUserData(token)
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		log.Println("Error fetching user data:", err)
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to format response", http.StatusInternalServerError)
		log.Println("Error marshaling response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// getAccessToken exchanges the authorization code for an access token.
func getAccessToken(code string) (string, error) {
	url := fmt.Sprintf("%s/access_token", oauthURL)
	data := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s&redirect_uri=%s", clientID, clientSecret, code, redirectURI)

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	values, err := urlValuesFromBody(string(body))
	if err != nil {
		return "", err
	}

	return values["access_token"], nil
}

// getUserData fetches user data from GitHub using the access token.
func getUserData(token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return nil, err
	}

	return userData, nil
}

// urlValuesFromBody parses a URL-encoded response body into a map.
func urlValuesFromBody(body string) (map[string]string, error) {
	values := make(map[string]string)
	for _, pair := range strings.Split(body, "&") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid body format")
		}
		values[kv[0]] = kv[1]
	}
	return values, nil
}
