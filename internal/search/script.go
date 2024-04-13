package search

import (
	"errors"
	"fmt"
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

func RunIndexTracks() error {
	esClient, err := SetupElasticClient()
	if err != nil {
		return err
	}

	if err := createTracksIndex(esClient, TrackIndex, "./assets/tracks.json"); err != nil {
		return err
	}

	return nil
}

func SetupElasticClient() (*elasticsearch.Client, error) {
	esURL := os.Getenv("ENV_ES_URL")
	if esURL == "" {
		return nil, errors.New("elasticsearch url cannot be empty")
	}

	// Use a third-party package for implementing the backoff function
	retryBackoff := backoff.NewExponentialBackOff()

	esConfig := elasticsearch.Config{
		Addresses:     []string{esURL},
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
