package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/jgsheppa/music_league_playlists/internal/playlists"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	esURL := os.Getenv("ENV_ES_URL")
	if esURL == "" {
		log.Fatal("elasticsearch url cannot be empty")
	}

	esConfig := elasticsearch.Config{
		Addresses: []string{esURL},
	}

	_, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("successfully connected to elastic node at the following url: %s \n", esURL)

	newPlaylist, err := playlists.NewPlaylist()
	if err != nil {
		log.Fatalf("could not initialize playlist: %e", err)
	}

	playlists, err := newPlaylist.GetPlaylistsForMulitpleLeagues()
	if err != nil {
		log.Fatalf("could not get playlists: %e", err)
	}

	jsonBytes, err := json.Marshal(playlists)
	if err != nil {
		log.Fatalf("could not marshal playlists into json: %e", err)
	}

	err = os.WriteFile("playlists.json", jsonBytes, 0644)
	if err != nil {
		log.Fatalf("could not create playlists.json: %e", err)
	}
}
