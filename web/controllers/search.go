package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jgsheppa/music_league_playlists/internal/search"
	"github.com/labstack/echo/v4"
)

func NewSearch(sc *search.ElasticSearch) *Search {
	return &Search{
		sc: sc,
	}
}

type Search struct {
	sc *search.ElasticSearch
}

func (s *Search) Result(c echo.Context) error {
	name := c.QueryParam("name")

	s.sc.WithIndex(search.TrackIndex)
	s.sc.WithQuery(name)
	res, err := s.sc.SearchField()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("could not decode json: %e", err)
	}

	var foundPlaylists search.TrackSearchResponse
	err = json.Unmarshal(body, &foundPlaylists)
	if err != nil {
		log.Fatalf("could not unmarshall json: %e", err)
	}

	return c.Render(http.StatusOK, "track", foundPlaylists.Hits)
}

func (s *Search) Home(c echo.Context) error {
	query := c.QueryParam("query")
	var foundPlaylists search.TrackSearchResponse
	fmt.Println(foundPlaylists)

	if len(query) > 0 {
		s.sc.WithIndex(search.TrackIndex)
		s.sc.WithQuery(query)
		res, err := s.sc.SearchField()
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("could not decode json: %e", err)
		}

		err = json.Unmarshal(body, &foundPlaylists)
		if err != nil {
			log.Fatalf("could not unmarshall json: %e", err)
		}
	}

	return c.Render(http.StatusOK, "home", foundPlaylists)
}
