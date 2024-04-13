package playlists_test

import (
	"errors"
	"os"
	"testing"

	"github.com/jgsheppa/music_league_playlists/internal/playlists"
)

func TestRemoveExistingFile(t *testing.T) {
	t.Parallel()

	filepath := "./testdata/test_playlist.txt"

	if err := playlists.RemoveExistingFile(filepath); err != nil {
		t.Fatalf("could not remove test file")
	}

	if _, err := os.Stat(filepath); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("found test playlist when it should not exist: %e", err)
	}

	if err := os.WriteFile(filepath, []byte("Everybody dance now"), 0644); err != nil {
		t.Errorf("could not create playlist: %e", err)
	}
}
