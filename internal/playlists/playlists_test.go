package playlists_test

import (
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
