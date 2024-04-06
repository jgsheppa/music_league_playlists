package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	humanize "github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/jgsheppa/music_league_playlists/internal/playlists"
)

const (
	PlaylistIndex = "playlist"
	SongIndex     = "song"
)

type ElasticSearch struct {
	client      *elasticsearch.Client
	bulkIndexer esutil.BulkIndexer
	err         error
	index       string
	filepath    string
}

func NewSearchClient(client *elasticsearch.Client) *ElasticSearch {
	return &ElasticSearch{client: client, err: nil, bulkIndexer: nil}
}

func (es *ElasticSearch) WithIndex(index string) *ElasticSearch {
	es.index = index
	return es
}

func (es *ElasticSearch) WithFile(filepath string) *ElasticSearch {
	es.filepath = filepath
	return es
}

func (es *ElasticSearch) WithBulkIndexer() *ElasticSearch {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         es.index,         // The default index name
		Client:        es.client,        // The Elasticsearch client
		NumWorkers:    runtime.NumCPU(), // The number of worker goroutines
		FlushBytes:    int(5e+6),        // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	})
	if err != nil {
		es.err = err
		return es
	}

	es.bulkIndexer = bi
	return es
}

func (es *ElasticSearch) indexDataInBulk(data []byte, countSuccessful uint64) {
	err := es.bulkIndexer.Add(
		context.Background(),
		esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewReader(data),
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				atomic.AddUint64(&countSuccessful, 1)
			},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: %s", err)
				} else {
					log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
				}
			},
		},
	)
	if err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}

func (es *ElasticSearch) displayBulkIndexerStats(start time.Time) {
	biStats := es.bulkIndexer.Stats()
	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Sucessfully indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}
}

func (es *ElasticSearch) CreateIndex(index string) error {
	_, err := es.client.Indices.Create(index)
	if err != nil {
		return err
	}
	return nil
}

func (es *ElasticSearch) RemoveIndex(index string) error {
	_, err := es.client.Indices.Delete([]string{index})
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

	for i, playlist := range jsonData {
		fmt.Println(i)
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

func Run() error {
	esURL := os.Getenv("ENV_ES_URL")
	if esURL == "" {
		return errors.New("elasticsearch url cannot be empty")
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
		return err
	}
	fmt.Printf("successfully connected to elastic node at the following url: %s \n", esURL)

	if err := createIndex(esClient, PlaylistIndex, "./assets/playlists.json"); err != nil {
		return err
	}

	if err := createIndex(esClient, SongIndex, "./assets/songs.json"); err != nil {
		return err
	}

	return nil

}

func createIndex(esClient *elasticsearch.Client, index, filepath string) error {
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
