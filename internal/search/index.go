package search

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jgsheppa/music_league_playlists/internal/playlists"
	"github.com/jgsheppa/music_league_playlists/internal/tracks"
)

const (
	PlaylistIndex = "playlist"
	TrackIndex    = "track"
)

func (es *ElasticSearch) CreateIndex(index string) error {
	_, err := es.client.Indices.Create(index)
	if err != nil {
		return err
	}
	return nil
}

func (es *ElasticSearch) RemoveIndex(index string) error {
	_, err := es.client.Indices.Exists([]string{index})
	if err != nil {
		log.Printf("could not check if index exists. status code")
		return nil
	}

	_, err = es.client.Indices.Delete([]string{index})
	if err != nil {
		return err
	}
	return nil
}

func (es *ElasticSearch) IndexPlaylists() error {
	var jsonData []playlists.SpotifyPlaylist
	data, err := os.ReadFile(es.filepath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	var countSuccessful uint64
	start := time.Now().UTC()

	for _, playlist := range jsonData {
		data, err := json.Marshal(playlist)
		if err != nil {
			return err
		}
		es.indexDataInBulk(data, countSuccessful)
	}

	if err := es.bulkIndexer.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	es.displayBulkIndexerStats(start)

	return nil
}

func (es *ElasticSearch) IndexTracks() error {
	var jsonData tracks.Items
	data, err := os.ReadFile(es.filepath)
	if err != nil {
		log.Printf("could not read file: %e", err)
		return err
	}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Printf("could not unmarshal json: %e", err)
		return err
	}

	var countSuccessful uint64
	start := time.Now().UTC()

	for _, playlist := range jsonData {
		data, err := json.Marshal(playlist)
		if err != nil {
			log.Printf("could not marshal json: %e", err)

			return err
		}
		es.indexDataInBulk(data, countSuccessful)
	}

	if err := es.bulkIndexer.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	es.displayBulkIndexerStats(start)

	return nil
}

func createPlaylistIndex(esClient *elasticsearch.Client, index, filepath string) error {
	client := NewSearchClient(esClient)
	client.WithIndex(index)
	client.WithBulkIndexer()
	client.WithFile(filepath)

	if err := client.RemoveIndex(client.index); err != nil {
		return err
	}

	if err := client.CreateIndex(client.index); err != nil {
		return err
	}

	if err := client.IndexPlaylists(); err != nil {
		return err
	}
	return nil
}

func createTracksIndex(esClient *elasticsearch.Client, index, filepath string) error {
	client := NewSearchClient(esClient)
	client.WithIndex(index)
	client.WithBulkIndexer()
	client.WithFile(filepath)

	if err := client.RemoveIndex(client.index); err != nil {
		log.Printf("could not remove index %e", err)
		return err
	}

	if err := client.CreateIndex(client.index); err != nil {
		log.Printf("could not create index %e", err)

		return err
	}

	if err := client.IndexTracks(); err != nil {
		log.Printf("could not index tracks %e", err)

		return err
	}
	return nil
}
