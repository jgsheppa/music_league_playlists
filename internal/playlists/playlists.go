package playlists

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jgsheppa/music_league_playlists/internal/spotify"
)

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

type PlaylistsByLeague = map[string][]SpotifyPlaylist

func MergeLeagueData(directory, filename string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	var mergedList MusicLeaguePlaylists
	for _, file := range files {
		if strings.Contains(file.Name(), ".json") {
			data, err := os.ReadFile(fmt.Sprintf("%s/%s", directory, file.Name()))
			if err != nil {
				return err
			}
			var list MusicLeaguePlaylists
			if err := json.Unmarshal(data, &list); err != nil {
				return err
			}

			mergedList = append(mergedList, list...)
		}
	}

	jsonBytes, err := json.Marshal(mergedList)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s", directory, filename)

	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		return nil
	}

	err = os.WriteFile(filePath, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *Playlist) GetPlaylists(filepath string) ([]SpotifyPlaylist, error) {
	start := time.Now()
	ids, err := ReadTestIDs(filepath)
	if err != nil {
		return []SpotifyPlaylist{}, err
	}

	res := NewWorker()
	jobs := make(chan string, len(ids))
	results := make(chan []SpotifyPlaylist, len(ids))
	var wg sync.WaitGroup

	numWorkers := 20
	for w := 0; w < numWorkers; w++ {
		go res.worker(p.getSpotifyPlaylist, jobs, results, &wg)
	}
	wg.Add(len(ids))

	for _, id := range ids {
		jobs <- id
	}

	close(jobs)
	wg.Wait()

	finish := time.Since(start).Seconds()
	fmt.Println("successfully fetched all playlists")
	fmt.Printf("finish in %f seconds \n", finish)
	return res.playlists, nil
}

type Playlist struct {
	spotify.SpotifyClient
}

func NewPlaylist() (*Playlist, error) {
	spotifyClient := spotify.NewSpotifyClient()

	_, err := spotifyClient.WithToken()
	if err != nil {
		return &Playlist{}, err
	}

	return &Playlist{SpotifyClient: *spotifyClient}, nil
}

func (p *Playlist) getSpotifyPlaylist(url string) (SpotifyPlaylist, error) {
	playlist, err := spotify.CreateSpotifyRequest[SpotifyPlaylist](url, p.Token)
	if err != nil {
		return SpotifyPlaylist{}, err
	}

	return playlist, nil
}
