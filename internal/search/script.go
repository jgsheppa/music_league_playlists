package search

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
)

func RunIndexPlaylists() error {
	esClient, err := SetupElasticClient()
	if err != nil {
		return err
	}

	if err := createPlaylistIndex(esClient, PlaylistIndex, "./assets/playlists.json"); err != nil {
		return err
	}

	return nil
}

func RunIndexTracks(esClient *elasticsearch.Client) error {
	if err := createTracksIndex(esClient, TrackIndex, "./web/assets/tracks.json"); err != nil {
		log.Printf("could not create tracks index %e", err)
		return err
	}

	return nil
}

func SetupElasticClient() (*elasticsearch.Client, error) {
	isProd := os.Getenv("IS_PROD")
	if isProd == "" {
		return nil, errors.New("environment cannot be empty")
	}

	apiKey := os.Getenv("ESC_API_KEY")
	if apiKey == "" && isProd == "true" {
		return nil, errors.New("elasticsearch api key cannot be empty")
	}

	cloudID := os.Getenv("ES_CLOUD_ID")
	if cloudID == "" && isProd == "true" {
		return nil, errors.New("elasticsearch cloud id cannot be empty")
	}

	// Use a third-party package for implementing the backoff function
	retryBackoff := backoff.NewExponentialBackOff()

	esConfig := elasticsearch.Config{
		CloudID:       cloudID,
		APIKey:        apiKey,
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	}

	if isProd == "false" {
		esConfig.Addresses = append(esConfig.Addresses, "http://localhost:9200")
		esConfig.APIKey = ""
		esConfig.CloudID = ""
	}

	esClient, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, err
	}
	fmt.Println("successfully connected to elastic node")

	return esClient, nil
}
