package main

import (
	"encoding/json"
	"io"
	"log"

	"github.com/jgsheppa/music_league_playlists/internal/search"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	esClient, err := search.SetupElasticClient()
	if err != nil {
		log.Fatal(err)
	}
	client := search.NewSearchClient(esClient)

	client.WithIndex(search.TrackIndex)
	client.WithQuery("Dio")
	res, err := client.SearchField()
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("could not decode json: %e", err)
	}

	var foundPlaylists search.TrackSearchResponse
	err = json.Unmarshal(body, &foundPlaylists)
	if err != nil {
		log.Fatalf("could not unmarshall json: %e", err)
	}
}
