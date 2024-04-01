package playlists

import (
	"errors"
	"fmt"
	"os"
)

func Run() error {
	if err := RemoveExistingPlaylist("./playlists.json"); err != nil {
		return err
	}

	if err := RemoveExistingPlaylist("./assets/music-league.json"); err != nil {
		return err
	}

	err := MergeLeagueData("assets", "music-league.json")
	if err != nil {
		return errors.Join(errors.New("could not merge music league playlists"), err)
	}

	playlistConfig, err := NewPlaylist()
	if err != nil {
		return errors.Join(errors.New("could not initialize playlist"), err)

	}

	playlists, err := playlistConfig.GetPlaylists(fmt.Sprintf("./assets/%s", "music-league.json"))
	if err != nil {
		return errors.Join(errors.New("could not get playlists"), err)
	}

	err = playlistConfig.CreatePlaylistJSON(playlists)
	if err != nil {
		return errors.Join(errors.New("could create merged playlist json file"), err)
	}
	return nil
}

func RemoveExistingPlaylist(filepath string) error {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err := os.Remove(filepath); err != nil {
		return err
	}

	return nil
}
