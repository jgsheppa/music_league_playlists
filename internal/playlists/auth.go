package playlists

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func getToken() (Token, error) {
	var token Token

	clientSecret := os.Getenv("SPOTIFY_SECRET")
	clientId := os.Getenv("SPOTIFY_CLIENT")

	data := strings.NewReader(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", clientId, clientSecret))

	req, err := http.NewRequest("POST", tokenURL, data)
	if err != nil {
		return Token{}, err
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&token)

	return token, nil
}
