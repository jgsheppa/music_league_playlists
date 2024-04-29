package main

import (
	"log"

	"github.com/jgsheppa/music_league_playlists/internal/search"
	"github.com/jgsheppa/music_league_playlists/web/controllers"
	"github.com/jgsheppa/music_league_playlists/web/views"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	esClient, err := search.SetupElasticClient()
	if err != nil {
		log.Fatal(err)
	}

	if err := search.RunIndexTracks(esClient); err != nil {
		log.Fatalf("could not index tracks: %e", err)
	}

	client := search.NewSearchClient(esClient)

	t := views.NewTemplate()

	searchController := controllers.NewSearch(client)

	e := echo.New()
	e.Renderer = t
	e.GET("/", searchController.Home)
	e.GET("/stats", searchController.Stats)
	e.Static("/static", "web/static")

	e.Logger.Fatal(e.Start(":8080"))
}
