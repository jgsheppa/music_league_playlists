package playlists

import (
	"errors"

	"github.com/jgsheppa/batbelt"
)

func Run() error {
	belt := batbelt.NewBatbelt()
	belt = belt.RemoveFile("./assets/playlists.json")
	if belt.Error() != nil {
		return belt.Error()
	}

	musicLeageFileName := "./assets/music-league.json"

	belt = belt.RemoveFile(musicLeageFileName)
	if belt.Error() != nil {
		return belt.Error()
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

	err = CreatePlaylistIDFile(musicLeageFileName, playlistIDsFilePath)
	if err != nil {
		return errors.Join(errors.New("could not create playlist id file"), err)
	}

	playlists, err := playlistClient.GetPlaylists(playlistIDsFilePath)
	if err != nil {
		return errors.Join(errors.New("could not get playlists"), err)
	}

	belt.CreateJSONFile(playlists, "web/assets/playlists.json")
	if belt.Error() != nil {
		return errors.Join(errors.New("could create merged playlist json file"), belt.Error())
	}

	return nil
}
