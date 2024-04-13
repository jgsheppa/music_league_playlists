package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jgsheppa/music_league_playlists/internal/playlists"
	"github.com/jgsheppa/music_league_playlists/internal/search"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	esClient, err := search.SetupElasticClient()
	if err != nil {
		log.Fatal(err)
	}
	client := search.NewSearchClient(esClient)
	client.WithIndex("playlist")
	client.WithQuery("tracks.items.track.name", "Surrender My Heart")
	res, err := client.SearchField()
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(res.Body)
	var result search.SearchResponse

	err = decoder.Decode(&result)
	if err != nil {
		log.Fatalf("could not decode json: %e", err)
	}

	b, err := json.Marshal(result.Hits.Hits)
	if err != nil {
		log.Fatalf("could not marshall json: %e", err)
	}

	var foundPlaylists []playlists.SpotifyPlaylist
	err = json.Unmarshal(b, &foundPlaylists)
	if err != nil {
		log.Fatalf("could not unmarshall json: %e", err)
	}
	fmt.Println(foundPlaylists)

}
