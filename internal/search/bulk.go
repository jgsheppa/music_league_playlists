package search

import (
	"bytes"
	"context"
	"log"
	"sync/atomic"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

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
		log.Printf("Unexpected error: %s", err)
	}
}

func (es *ElasticSearch) displayBulkIndexerStats(start time.Time) {
	biStats := es.bulkIndexer.Stats()
	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Printf(
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
