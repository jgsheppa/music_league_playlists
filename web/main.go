package main

import (
	"log"

	"github.com/jgsheppa/music_league_playlists/internal/playlists"
	"github.com/jgsheppa/music_league_playlists/internal/search"
	"github.com/jgsheppa/music_league_playlists/internal/tracks"
	"github.com/jgsheppa/music_league_playlists/web/controllers"
	"github.com/jgsheppa/music_league_playlists/web/views"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func init() {
	if err := playlists.Run(); err != nil {
		log.Fatalf("could not run playlist script to completion: %e", err)
	}

	if err := search.RunIndexPlaylists(); err != nil {
		log.Fatalf("could not run search script: %e", err)
	}

	if err := tracks.Run(); err != nil {
		log.Fatalf("could not run track command: %e", err)
	}

	if err := search.RunIndexTracks(); err != nil {
		log.Fatalf("could not run search script: %e", err)
	}
}
func main() {
	esClient, err := search.SetupElasticClient()
	if err != nil {
		log.Fatal(err)
	}
	client := search.NewSearchClient(esClient)

	t := views.NewTemplate()

	searchController := controllers.NewSearch(client)

	e := echo.New()
	e.Renderer = t
	e.GET("/", searchController.Home)
	e.Logger.Fatal(e.Start(":8080"))
}
