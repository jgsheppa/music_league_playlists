package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SpotifyClient struct {
	Token Token
}

const (
	PlaylistURL = "https://api.spotify.com/v1/playlists/"
)

func NewSpotifyClient() *SpotifyClient {
	return &SpotifyClient{}
}

func (sc *SpotifyClient) WithToken() (*SpotifyClient, error) {
	token, err := getToken()
	if err != nil {
		return &SpotifyClient{}, err
	}
	sc.Token = token
	return sc, nil
}

func CreateSpotifyRequest[T any](url string, token Token) (T, error) {
	var result T
	req, err := http.NewRequest("GET", url, &strings.Reader{})
	if err != nil {
		return result, err
	}

	client := &http.Client{}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}
