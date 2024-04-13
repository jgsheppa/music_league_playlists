package spotify

type SpotifyClient struct {
	Token Token
}

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
