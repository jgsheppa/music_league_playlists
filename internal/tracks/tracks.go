package tracks

type Tracks struct {
	Href     string `json:"href,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Next     string `json:"next,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	Previous string `json:"previous,omitempty"`
	Total    int    `json:"total,omitempty"`
	Items    []struct {
		AddedAt string `json:"added_at,omitempty"`
		AddedBy struct {
			ExternalUrls struct {
				Spotify string `json:"spotify,omitempty"`
			} `json:"external_urls,omitempty"`
			Followers struct {
				Href  string `json:"href,omitempty"`
				Total int    `json:"total,omitempty"`
			} `json:"followers,omitempty"`
			Href string `json:"href,omitempty"`
			ID   string `json:"id,omitempty"`
			Type string `json:"type,omitempty"`
			URI  string `json:"uri,omitempty"`
		} `json:"added_by,omitempty"`
		IsLocal bool `json:"is_local,omitempty"`
		Track   struct {
			Album struct {
				AlbumType        string   `json:"album_type,omitempty"`
				TotalTracks      int      `json:"total_tracks,omitempty"`
				AvailableMarkets []string `json:"available_markets,omitempty"`
				ExternalUrls     struct {
					Spotify string `json:"spotify,omitempty"`
				} `json:"external_urls,omitempty"`
				Href   string `json:"href,omitempty"`
				ID     string `json:"id,omitempty"`
				Images []struct {
					URL    string `json:"url,omitempty"`
					Height int    `json:"height,omitempty"`
					Width  int    `json:"width,omitempty"`
				} `json:"images,omitempty"`
				Name                 string `json:"name,omitempty"`
				ReleaseDate          string `json:"release_date,omitempty"`
				ReleaseDatePrecision string `json:"release_date_precision,omitempty"`
				Restrictions         struct {
					Reason string `json:"reason,omitempty"`
				} `json:"restrictions,omitempty"`
				Type    string `json:"type,omitempty"`
				URI     string `json:"uri,omitempty"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify,omitempty"`
					} `json:"external_urls,omitempty"`
					Href string `json:"href,omitempty"`
					ID   string `json:"id,omitempty"`
					Name string `json:"name,omitempty"`
					Type string `json:"type,omitempty"`
					URI  string `json:"uri,omitempty"`
				} `json:"artists,omitempty"`
			} `json:"album,omitempty"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify,omitempty"`
				} `json:"external_urls,omitempty"`
				Followers struct {
					Href  string `json:"href,omitempty"`
					Total int    `json:"total,omitempty"`
				} `json:"followers,omitempty"`
				Genres []string `json:"genres,omitempty"`
				Href   string   `json:"href,omitempty"`
				ID     string   `json:"id,omitempty"`
				Images []struct {
					URL    string `json:"url,omitempty"`
					Height int    `json:"height,omitempty"`
					Width  int    `json:"width,omitempty"`
				} `json:"images,omitempty"`
				Name       string `json:"name,omitempty"`
				Popularity int    `json:"popularity,omitempty"`
				Type       string `json:"type,omitempty"`
				URI        string `json:"uri,omitempty"`
			} `json:"artists,omitempty"`
			AvailableMarkets []string `json:"available_markets,omitempty"`
			DiscNumber       int      `json:"disc_number,omitempty"`
			DurationMs       int      `json:"duration_ms,omitempty"`
			Explicit         bool     `json:"explicit,omitempty"`
			ExternalIds      struct {
				Isrc string `json:"isrc,omitempty"`
				Ean  string `json:"ean,omitempty"`
				Upc  string `json:"upc,omitempty"`
			} `json:"external_ids,omitempty"`
			ExternalUrls struct {
				Spotify string `json:"spotify,omitempty"`
			} `json:"external_urls,omitempty"`
			Href       string `json:"href,omitempty"`
			ID         string `json:"id,omitempty"`
			IsPlayable bool   `json:"is_playable,omitempty"`
			LinkedFrom struct {
			} `json:"linked_from,omitempty"`
			Restrictions struct {
				Reason string `json:"reason,omitempty"`
			} `json:"restrictions,omitempty"`
			Name        string `json:"name,omitempty"`
			Popularity  int    `json:"popularity,omitempty"`
			PreviewURL  string `json:"preview_url,omitempty"`
			TrackNumber int    `json:"track_number,omitempty"`
			Type        string `json:"type,omitempty"`
			URI         string `json:"uri,omitempty"`
			IsLocal     bool   `json:"is_local,omitempty"`
		} `json:"track,omitempty"`
	} `json:"items,omitempty"`
}
