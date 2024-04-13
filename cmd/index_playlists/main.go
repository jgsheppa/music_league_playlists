package main

import (
	"log"

	"github.com/jgsheppa/music_league_playlists/internal/playlists"
	"github.com/jgsheppa/music_league_playlists/internal/search"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if err := playlists.Run(); err != nil {
		log.Fatalf("could not run playlist script to completion: %e", err)
	}

	if err := search.RunIndexPlaylists(); err != nil {
		log.Fatalf("could not run search script: %e", err)
	}
}
