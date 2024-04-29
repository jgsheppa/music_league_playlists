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

func (s *Search) Home(c echo.Context) error {
	query := c.QueryParam("query")
	var foundPlaylists search.TrackSearchResponse

	if len(query) > 0 {
		s.sc.WithIndex(search.TrackIndex)
		s.sc.WithQuery(query)
		res, err := s.sc.SearchField()
		if err != nil {
			log.Println(err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("could not decode json: %e \n", err)
		}

		err = json.Unmarshal(body, &foundPlaylists)
		if err != nil {
			log.Printf("could not unmarshall json: %e \n", err)
		}
	}

	return c.Render(http.StatusOK, "home", foundPlaylists)
}

func (s *Search) Stats(c echo.Context) error {
	s.sc.WithIndex(search.TrackIndex)
	res, err := s.sc.AggregateByTerm()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res)

	return c.JSON(http.StatusOK, res)
}
