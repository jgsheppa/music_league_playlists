package main

import (
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
	playlist, err := newPlaylist.GetPlaylists("assets/electricBoogaloo.json")
	if err != nil {
		log.Fatalf("could not get playlists: %e", err)
	}
	fmt.Println(playlist)
}
