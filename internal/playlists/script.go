package playlists

import (
	"errors"
	"os"
)

func Run() error {
	if err := RemoveExistingPlaylist("./assets/playlists.json"); err != nil {
		return err
	}

	if err := RemoveExistingPlaylist("./assets/music-league.json"); err != nil {
		return err
	}

	err := MergeLeagueData("assets", "music-league.json")
	if err != nil {
		return errors.Join(errors.New("could not merge music league playlists"), err)
	}

	playlistClient, err := NewPlaylist()
	if err != nil {
		return errors.Join(errors.New("could not initialize playlist"), err)
	}

	playlistIDsFilePath := "./assets/playlist-ids.txt"

	err = CreatePlaylistIDFile("./assets/music-league.json", playlistIDsFilePath)
	if err != nil {
		return errors.Join(errors.New("could not create playlist id file"), err)
	}

	playlists, err := playlistClient.GetPlaylists(playlistIDsFilePath)
	if err != nil {
		return errors.Join(errors.New("could not get playlists"), err)
	}

	err = playlistClient.CreatePlaylistJSON(playlists)
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
