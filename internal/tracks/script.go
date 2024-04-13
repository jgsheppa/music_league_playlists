package tracks

import (
	"errors"

	"github.com/jgsheppa/music_league_playlists/internal/playlists"
)

func Run() error {
	if err := playlists.RemoveExistingFile("./assets/tracks.json"); err != nil {
		return err
	}

	trackClient, err := NewTrack()
	if err != nil {
		return errors.Join(errors.New("could not create track client"), err)
	}

	playlistIDsFilePath := "./assets/playlist-ids.txt"

	tracks, err := trackClient.GetTracks(playlistIDsFilePath)
	if err != nil {
		return errors.Join(errors.New("could not get spotify tracks"), err)
	}

	err = trackClient.CreateTrackJSON(tracks)
	if err != nil {
		return errors.Join(errors.New("could not get spotify tracks"), err)
	}

	return nil
}
