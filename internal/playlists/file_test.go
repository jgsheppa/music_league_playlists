package playlists_test

import (
	"errors"
	"os"
	"testing"

	"github.com/jgsheppa/music_league_playlists/internal/playlists"
)

func TestCreatePlaylistIDFile(t *testing.T) {
	t.Parallel()

	idFilepath := "./testdata/test_playlist_ids.txt"
	playlistFilepath := "../../assets/electricBoogaloo.json"

	err := playlists.CreatePlaylistIDFile(playlistFilepath, idFilepath)
	if err != nil {
		t.Errorf("could not create id file: %e", err)
	}

	if _, err = os.Stat(idFilepath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("could not find created playlist: %e", err)
	}

	if err = os.Remove(idFilepath); err != nil {
		t.Errorf("could not remove playlist: %e", err)
	}
}

func TestReadTestIDs(t *testing.T) {
	t.Parallel()

	idFilepath := "./testdata/test_playlist_ids.txt"
	playlistFilepath := "../../assets/electricBoogaloo.json"

	err := playlists.CreatePlaylistIDFile(playlistFilepath, idFilepath)
	if err != nil {
		t.Errorf("could not create id file: %e", err)
	}

	ids, err := playlists.ReadTestIDs(idFilepath)
	if err != nil {
		t.Errorf("could not read id file: %e", err)
	}

	if len(ids) == 0 {
		t.Errorf("no ids in test file: %e", err)
	}

	want := "4TdmvoQKLOdSPwp7o1mCFB"
	got := ids[0]

	if want != got {
		t.Errorf("want %v not equal to got %v", want, got)
	}
}
