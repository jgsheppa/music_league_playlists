package search

import (
	"fmt"
	"runtime"
	"strings"
	"time"

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
	value string
}

func NewSearchClient(client *elasticsearch.Client) *ElasticSearch {
	return &ElasticSearch{client: client, err: nil, bulkIndexer: nil}
}

func (es *ElasticSearch) WithIndex(index string) *ElasticSearch {
	es.index = index
	return es
}

func (es *ElasticSearch) WithQuery(value string) *ElasticSearch {
	es.query = Query{
		value,
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

func (es *ElasticSearch) SearchField() (*esapi.Response, error) {
	multiQuery := fmt.Sprintf(`{
	  "query": {
	    "multi_match": {
	      "query": "%s",
	      "type": "most_fields",
	      "fields": [
	        "track.artists.name" ,
					"track.album.name",
	        "track.name"
	      ],
				"fuzziness": "Auto",
				"prefix_length": "2"
	    }
	  }
	}`, es.query.value)
	res, err := es.client.Search(es.client.Search.WithIndex(es.index), es.client.Search.WithBody(strings.NewReader(multiQuery)))

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

type TrackSearchResponse struct {
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
		MaxScore any         `json:"max_score,omitempty"`
		Hits     []FoundItem `json:"hits"`
	} `json:"hits,omitempty"`
}

type FoundItem struct {
	ID     string  `json:"_id"`
	Index  string  `json:"_index"`
	Score  float64 `json:"_score"`
	Source struct {
		AddedAt time.Time `json:"added_at"`
		AddedBy struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
			} `json:"followers"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"added_by"`
		Track struct {
			Album struct {
				AlbumType string `json:"album_type"`
				Artists   []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
				AvailableMarkets []string `json:"available_markets"`
				ExternalUrls     struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				ID     string `json:"id"`
				Images []struct {
					Height int    `json:"height"`
					URL    string `json:"url"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				Restrictions         struct {
				} `json:"restrictions"`
				TotalTracks int    `json:"total_tracks"`
				Type        string `json:"type"`
				URI         string `json:"uri"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Followers struct {
				} `json:"followers"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			DiscNumber       int      `json:"disc_number"`
			DurationMs       int      `json:"duration_ms"`
			ExternalIds      struct {
				Isrc string `json:"isrc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href       string `json:"href"`
			ID         string `json:"id"`
			LinkedFrom struct {
			} `json:"linked_from"`
			Name         string `json:"name"`
			Popularity   int    `json:"popularity"`
			PreviewURL   string `json:"preview_url"`
			Restrictions struct {
			} `json:"restrictions"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
		} `json:"track"`
	} `json:"_source"`
	Type string `json:"_type"`
}
