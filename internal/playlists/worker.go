package playlists

import (
	"fmt"
	"sync"

	"github.com/jgsheppa/music_league_playlists/internal/spotify"
)

type Result struct {
	playlists []SpotifyPlaylist
	err       error
}

func NewWorker() *Result {
	return &Result{}
}

func (r *Result) worker(requester func(id string) (SpotifyPlaylist, error), jobs <-chan string, results chan<- []SpotifyPlaylist, wg *sync.WaitGroup) {
	r.playlists = make([]SpotifyPlaylist, len(jobs))
	for job := range jobs {
		res, err := requester(fmt.Sprintf("%s%s", spotify.PlaylistURL, job))
		if err != nil {
			r.err = err
		}

		if res.ExternalUrls.Spotify == "" {
			fmt.Printf("empty playlist id: %s", job)
		}

		r.playlists = append(r.playlists, res)
		results <- r.playlists
		wg.Done()
	}
}
