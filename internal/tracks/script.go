package tracks

import (
	"errors"

	"github.com/jgsheppa/batbelt"
)

func Run() error {
	belt := batbelt.NewBatbelt()
	belt = belt.RemoveFile("./web/assets/tracks.json")
	if belt.Error() != nil {
		return belt.Error()
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

	belt.CreateJSONFile(tracks, "web/assets/tracks.json")
	if belt.Error() != nil {
		return errors.Join(errors.New("could not get spotify tracks"), belt.Error())
	}

	return nil
}
