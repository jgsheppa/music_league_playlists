package tracks

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jgsheppa/music_league_playlists/internal/playlists"
	"github.com/jgsheppa/music_league_playlists/internal/spotify"
)

type SpotifyTracks struct {
	Href     string `json:"href,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Next     string `json:"next,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	Previous string `json:"previous,omitempty"`
	Total    int    `json:"total,omitempty"`
	Items    Items  `json:"items,omitempty"`
}

type Items = []Item

type Item struct {
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
}

type Track struct {
	spotify.SpotifyClient
}

func NewTrack() (*Track, error) {
	spotifyClient := spotify.NewSpotifyClient()

	_, err := spotifyClient.WithToken()
	if err != nil {
		return &Track{}, err
	}

	return &Track{SpotifyClient: *spotifyClient}, nil
}

func (t *Track) getSpotifyTracks(url string) (Items, error) {
	playlist, err := spotify.CreateSpotifyRequest[SpotifyTracks](url, t.SpotifyClient.Token)
	if err != nil {
		return Items{}, err
	}

	return playlist.Items, nil
}

func (t *Track) GetTracks(filepath string) (Items, error) {
	if t.Token.AccessToken == "" {
		return Items{}, errors.New("no spotify access token found")
	}

	start := time.Now()
	ids, err := playlists.ReadTestIDs(filepath)
	if err != nil {
		return Items{}, err
	}

	res := NewWorker()
	jobs := make(chan string, len(ids))
	results := make(chan Items, len(ids))
	var wg sync.WaitGroup

	numWorkers := 20
	for w := 0; w < numWorkers; w++ {
		go res.worker(t.getSpotifyTracks, jobs, results, &wg)
	}
	wg.Add(len(ids))

	for _, id := range ids {
		jobs <- id
	}

	close(jobs)
	wg.Wait()

	finish := time.Since(start).Seconds()
	fmt.Println("successfully fetched all tracks")
	fmt.Printf("finish in %f seconds \n", finish)
	return res.tracks, nil
}
