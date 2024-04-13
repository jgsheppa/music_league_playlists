package search

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type ElasticSearch struct {
	client      *elasticsearch.Client
	bulkIndexer esutil.BulkIndexer
	err         error
	index       string
	filepath    string
	query       Query
}

type Query struct {
	field string
	value string
}

func NewSearchClient(client *elasticsearch.Client) *ElasticSearch {
	return &ElasticSearch{client: client, err: nil, bulkIndexer: nil}
}

func (es *ElasticSearch) WithIndex(index string) *ElasticSearch {
	es.index = index
	return es
}

func (es *ElasticSearch) WithQuery(field, value string) *ElasticSearch {
	es.query = Query{
		field, value,
	}
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

func (es *ElasticSearch) SearchField() (*esapi.Response, error) {

	query := fmt.Sprintf(`{ "query": {"bool": {"must": [{ "match": {"%s" : {"query": "%s"} } }] } } }`, es.query.field, es.query.value)
	res, err := es.client.Search(es.client.Search.WithIndex(es.index), es.client.Search.WithBody(strings.NewReader(query)))

	if err != nil {
		return nil, err
	}
	return res, nil
}

type SearchResponse struct {
	Took     int  `json:"took,omitempty"`
	TimedOut bool `json:"timed_out,omitempty"`
	Shards   struct {
		Total      int `json:"total,omitempty"`
		Successful int `json:"successful,omitempty"`
		Skipped    int `json:"skipped,omitempty"`
		Failed     int `json:"failed,omitempty"`
	} `json:"_shards,omitempty"`
	Hits struct {
		Total struct {
			Value    int    `json:"value,omitempty"`
			Relation string `json:"relation,omitempty"`
		} `json:"total,omitempty"`
		MaxScore any   `json:"max_score,omitempty"`
		Hits     []any `json:"hits"`
	} `json:"hits,omitempty"`
}
