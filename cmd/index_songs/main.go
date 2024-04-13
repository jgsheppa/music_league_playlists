package main

import (
	"log"

	"github.com/jgsheppa/music_league_playlists/internal/tracks"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if err := tracks.Run(); err != nil {
		log.Fatalf("could not run track command: %e", err)
	}
}
