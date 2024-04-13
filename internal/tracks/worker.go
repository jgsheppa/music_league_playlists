package tracks

import (
	"fmt"
	"sync"

	"github.com/jgsheppa/music_league_playlists/internal/spotify"
)

type Result struct {
	tracks Items
	err    error
}

func NewWorker() *Result {
	return &Result{}
}

func (r *Result) worker(requester func(id string) (Items, error), jobs <-chan string, results chan<- Items, wg *sync.WaitGroup) {
	r.tracks = make(Items, len(jobs))
	for job := range jobs {
		res, err := requester(fmt.Sprintf("%s%s", spotify.PlaylistURL, job))
		if err != nil {
			r.err = err
		}

		r.tracks = append(r.tracks, res...)
		results <- r.tracks
		wg.Done()
	}
}
