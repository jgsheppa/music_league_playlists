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
	if err := createTracksIndex(esClient, TrackIndex, "./assets/tracks.json"); err != nil {
		log.Printf("could not create tracks index %e", err)
		return err
	}

	return nil
}

func SetupElasticClient() (*elasticsearch.Client, error) {
	esURL := os.Getenv("ES_URL")
	if esURL == "" {
		return nil, errors.New("elasticsearch url cannot be empty")
	}

	apiKey := os.Getenv("ESC_API_KEY")
	if esURL == "" {
		return nil, errors.New("elasticsearch api key cannot be empty")
	}

	cloudID := os.Getenv("ES_CLOUD_ID")
	if esURL == "" {
		return nil, errors.New("elasticsearch cloud id cannot be empty")
	}

	// Use a third-party package for implementing the backoff function
	retryBackoff := backoff.NewExponentialBackOff()

	esConfig := elasticsearch.Config{
		// Addresses:     []string{esURL},
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

	esClient, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, err
	}
	fmt.Printf("successfully connected to elastic node at the following url: %s \n", esURL)

	return esClient, nil
}
