package playlists

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strings"
)

func GetPlaylistIDs(playlists []byte) ([]string, error) {
	var list MusicLeaguePlaylists
	if err := json.Unmarshal(playlists, &list); err != nil {
		return nil, err
	}

	var ids []string
	for _, playlist := range list {
		id := strings.Split(playlist.PlaylistURL, "playlist/")

		if len(id) == 2 {
			ids = append(ids, id[1])
		}
	}

	return ids, nil
}

func CreatePlaylistIDFile(playlistFilepath, idFilepath string) error {
	file, err := os.ReadFile(playlistFilepath)
	if err != nil {
		return err
	}

	ids, err := GetPlaylistIDs(file)
	if err != nil {
		return err
	}

	err = os.WriteFile(idFilepath, []byte(strings.Join(ids, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadTestIDs(filepath string) ([]string, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var ids []string
	scanner := bufio.NewScanner(bytes.NewReader(file))
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func RemoveExistingFile(filepath string) error {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err := os.Remove(filepath); err != nil {
		return err
	}

	return nil
}
