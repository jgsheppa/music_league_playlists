package playlists

import (
	"errors"
)

func Run() error {
	if err := RemoveExistingFile("./assets/playlists.json"); err != nil {
		return err
	}

	if err := RemoveExistingFile("./assets/music-league.json"); err != nil {
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
