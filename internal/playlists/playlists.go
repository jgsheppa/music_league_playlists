package playlists

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

//go:embed assets/*
var content embed.FS

type MusicLeaguePlaylist struct {
	ID                  string    `json:"id,omitempty"`
	Name                string    `json:"name,omitempty"`
	Completed           time.Time `json:"completed,omitempty"`
	Description         string    `json:"description,omitempty"`
	DownvotesPerUser    int       `json:"downvotesPerUser,omitempty"`
	HighStakes          bool      `json:"highStakes,omitempty"`
	LeagueID            string    `json:"leagueId,omitempty"`
	MaxDownvotesPerSong int       `json:"maxDownvotesPerSong,omitempty"`
	MaxUpvotesPerSong   int       `json:"maxUpvotesPerSong,omitempty"`
	PlaylistURL         string    `json:"playlistUrl,omitempty"`
	Sequence            int       `json:"sequence,omitempty"`
	SongsPerUser        int       `json:"songsPerUser,omitempty"`
	StartDate           time.Time `json:"startDate,omitempty"`
	Status              string    `json:"status,omitempty"`
	SubmissionsDue      time.Time `json:"submissionsDue,omitempty"`
	UpvotesPerUser      int       `json:"upvotesPerUser,omitempty"`
	VotesDue            time.Time `json:"votesDue,omitempty"`
	TemplateID          string    `json:"templateId,omitempty"`
}

type MusicLeaguePlaylists = []MusicLeaguePlaylist

type SpotifyPlaylist struct {
	Collaborative bool   `json:"collaborative,omitempty"`
	Description   string `json:"description,omitempty"`
	ExternalUrls  struct {
		Spotify string `json:"spotify,omitempty"`
	} `json:"external_urls,omitempty"`
	Followers struct {
		Href  string `json:"href,omitempty"`
		Total int    `json:"total,omitempty"`
	} `json:"followers,omitempty"`
	Href   string `json:"href,omitempty"`
	ID     string `json:"id,omitempty"`
	Images []struct {
		URL    string `json:"url,omitempty"`
		Height int    `json:"height,omitempty"`
		Width  int    `json:"width,omitempty"`
	} `json:"images,omitempty"`
	Name  string `json:"name,omitempty"`
	Owner struct {
		ExternalUrls struct {
			Spotify string `json:"spotify,omitempty"`
		} `json:"external_urls,omitempty"`
		Followers struct {
			Href  string `json:"href,omitempty"`
			Total int    `json:"total,omitempty"`
		} `json:"followers,omitempty"`
		Href        string `json:"href,omitempty"`
		ID          string `json:"id,omitempty"`
		Type        string `json:"type,omitempty"`
		URI         string `json:"uri,omitempty"`
		DisplayName string `json:"display_name,omitempty"`
	} `json:"owner,omitempty"`
	Public     bool   `json:"public,omitempty"`
	SnapshotID string `json:"snapshot_id,omitempty"`
	Tracks     struct {
		Href     string `json:"href,omitempty"`
		Limit    int    `json:"limit,omitempty"`
		Next     string `json:"next,omitempty"`
		Offset   int    `json:"offset,omitempty"`
		Previous string `json:"previous,omitempty"`
		Total    int    `json:"total,omitempty"`
		Items    []struct {
			AddedAt string `json:"added_at,omitempty"`
			AddedBy struct {
				ExternalUrls struct {
					Spotify string `json:"spotify,omitempty"`
				} `json:"external_urls,omitempty"`
				Followers struct {
					Href  string `json:"href,omitempty"`
					Total int    `json:"total,omitempty"`
				} `json:"followers,omitempty"`
				Href string `json:"href,omitempty"`
				ID   string `json:"id,omitempty"`
				Type string `json:"type,omitempty"`
				URI  string `json:"uri,omitempty"`
			} `json:"added_by,omitempty"`
			IsLocal bool `json:"is_local,omitempty"`
			Track   struct {
				Album struct {
					AlbumType        string   `json:"album_type,omitempty"`
					TotalTracks      int      `json:"total_tracks,omitempty"`
					AvailableMarkets []string `json:"available_markets,omitempty"`
					ExternalUrls     struct {
						Spotify string `json:"spotify,omitempty"`
					} `json:"external_urls,omitempty"`
					Href   string `json:"href,omitempty"`
					ID     string `json:"id,omitempty"`
					Images []struct {
						URL    string `json:"url,omitempty"`
						Height int    `json:"height,omitempty"`
						Width  int    `json:"width,omitempty"`
					} `json:"images,omitempty"`
					Name                 string `json:"name,omitempty"`
					ReleaseDate          string `json:"release_date,omitempty"`
					ReleaseDatePrecision string `json:"release_date_precision,omitempty"`
					Restrictions         struct {
						Reason string `json:"reason,omitempty"`
					} `json:"restrictions,omitempty"`
					Type    string `json:"type,omitempty"`
					URI     string `json:"uri,omitempty"`
					Artists []struct {
						ExternalUrls struct {
							Spotify string `json:"spotify,omitempty"`
						} `json:"external_urls,omitempty"`
						Href string `json:"href,omitempty"`
						ID   string `json:"id,omitempty"`
						Name string `json:"name,omitempty"`
						Type string `json:"type,omitempty"`
						URI  string `json:"uri,omitempty"`
					} `json:"artists,omitempty"`
				} `json:"album,omitempty"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify,omitempty"`
					} `json:"external_urls,omitempty"`
					Followers struct {
						Href  string `json:"href,omitempty"`
						Total int    `json:"total,omitempty"`
					} `json:"followers,omitempty"`
					Genres []string `json:"genres,omitempty"`
					Href   string   `json:"href,omitempty"`
					ID     string   `json:"id,omitempty"`
					Images []struct {
						URL    string `json:"url,omitempty"`
						Height int    `json:"height,omitempty"`
						Width  int    `json:"width,omitempty"`
					} `json:"images,omitempty"`
					Name       string `json:"name,omitempty"`
					Popularity int    `json:"popularity,omitempty"`
					Type       string `json:"type,omitempty"`
					URI        string `json:"uri,omitempty"`
				} `json:"artists,omitempty"`
				AvailableMarkets []string `json:"available_markets,omitempty"`
				DiscNumber       int      `json:"disc_number,omitempty"`
				DurationMs       int      `json:"duration_ms,omitempty"`
				Explicit         bool     `json:"explicit,omitempty"`
				ExternalIds      struct {
					Isrc string `json:"isrc,omitempty"`
					Ean  string `json:"ean,omitempty"`
					Upc  string `json:"upc,omitempty"`
				} `json:"external_ids,omitempty"`
				ExternalUrls struct {
					Spotify string `json:"spotify,omitempty"`
				} `json:"external_urls,omitempty"`
				Href       string `json:"href,omitempty"`
				ID         string `json:"id,omitempty"`
				IsPlayable bool   `json:"is_playable,omitempty"`
				LinkedFrom struct {
				} `json:"linked_from,omitempty"`
				Restrictions struct {
					Reason string `json:"reason,omitempty"`
				} `json:"restrictions,omitempty"`
				Name        string `json:"name,omitempty"`
				Popularity  int    `json:"popularity,omitempty"`
				PreviewURL  string `json:"preview_url,omitempty"`
				TrackNumber int    `json:"track_number,omitempty"`
				Type        string `json:"type,omitempty"`
				URI         string `json:"uri,omitempty"`
				IsLocal     bool   `json:"is_local,omitempty"`
			} `json:"track,omitempty"`
		} `json:"items,omitempty"`
	} `json:"tracks,omitempty"`
	Type string `json:"type,omitempty"`
	URI  string `json:"uri,omitempty"`
}

func GetPlaylistIDs(playlists []byte) ([]string, error) {
	var list MusicLeaguePlaylists
	if err := json.Unmarshal(playlists, &list); err != nil {
		return nil, err
	}

	var ids []string
	for _, playlist := range list {
		id := strings.Split(playlist.PlaylistURL, "playlist/")

		if len(id) == 2 {
			ids = append(ids, id[1])
		}
	}

	return ids, nil
}

const (
	SpotifyURL = "https://api.spotify.com/v1/playlists/"
	tokenURL   = "https://accounts.spotify.com/api/token"
)

type PlaylistsByLeague = map[string][]SpotifyPlaylist

func (p *Playlist) GetPlaylistsForMulitpleLeagues() ([]SpotifyPlaylist, error) {
	var fullList []SpotifyPlaylist

	files, err := content.ReadDir("assets")
	if err != nil {
		return []SpotifyPlaylist{}, err
	}

	var wg sync.WaitGroup
	for _, file := range files {
		fmt.Printf("fetching playlists for file: %s \n", file.Name())
		wg.Add(1)
		var playlists []SpotifyPlaylist

		go func() {
			defer wg.Done()
			playlists, err = p.GetPlaylists(fmt.Sprintf("assets/%s", file.Name()))
		}()
		wg.Wait()
		if err != nil {
			return []SpotifyPlaylist{}, err
		}
		fullList = append(fullList, playlists...)
	}
	fmt.Println("successfully fetched all playlists")
	return fullList, nil
}

func (p *Playlist) GetPlaylists(filepath string) ([]SpotifyPlaylist, error) {
	file, err := content.ReadFile(filepath)
	if err != nil {
		return []SpotifyPlaylist{}, err
	}

	ids, err := GetPlaylistIDs(file)
	if err != nil {
		return []SpotifyPlaylist{}, err
	}

	var wg sync.WaitGroup
	var playlists []SpotifyPlaylist
	for _, id := range ids {
		wg.Add(1)
		var playlist SpotifyPlaylist

		go func() {
			defer wg.Done()
			playlist, err = p.getPlaylist(id)
		}()
		wg.Wait()
		if err != nil {
			return []SpotifyPlaylist{}, err
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

type Playlist struct {
	token Token
}

type Token struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

func NewPlaylist() (*Playlist, error) {
	token, err := getToken()
	if err != nil {
		return &Playlist{}, err
	}

	return &Playlist{token}, nil
}

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

func (p *Playlist) getPlaylist(id string) (SpotifyPlaylist, error) {
	var playlist SpotifyPlaylist
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", SpotifyURL, id), &strings.Reader{})
	if err != nil {
		return SpotifyPlaylist{}, err
	}
	client := &http.Client{}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		return SpotifyPlaylist{}, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&playlist)

	return playlist, nil
}
