package playlists_test

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/jgsheppa/music_league_playlists/internal/playlists"
)

func TestGetPlaylistIDs(t *testing.T) {
	t.Parallel()

	file, err := os.ReadFile("../../assets/electricBoogaloo.json")
	if err != nil {
		t.Errorf("could not open .json file: %e", err)
	}

	got, err := playlists.GetPlaylistIDs(file)
	if err != nil {
		t.Errorf("could not get playlist ids: %e", err)
	}

	want := "4TdmvoQKLOdSPwp7o1mCFB"

	if !slices.Contains(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	want = "3VnnViGYtMU4m01psjtnoa"
	if !slices.Contains(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMergeLeagueData(t *testing.T) {
	t.Parallel()

	directory := "./testdata"
	filename := "test.json"
	filepath := fmt.Sprintf("%s/%s", directory, filename)

	err := playlists.MergeLeagueData(directory, filename)
	if err != nil {
		t.Errorf("could not merge playlists: %e", err)
	}

	if _, err = os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("could not find created playlist: %e", err)
	}

	if err = os.Remove(filepath); err != nil {
		t.Errorf("could not remove playlist: %e", err)
	}

}
