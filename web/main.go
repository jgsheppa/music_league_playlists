package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/jgsheppa/music_league_playlists/internal/search"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	esClient, err := search.SetupElasticClient()
	if err != nil {
		log.Fatal(err)
	}
	client := search.NewSearchClient(esClient)

	t := &Template{
		templates: template.Must(template.ParseGlob("web/views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.GET("/tracks", func(c echo.Context) error {
		name := c.QueryParam("name")

		client.WithIndex(search.TrackIndex)
		client.WithQuery(name)
		res, err := client.SearchField()
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
	})
	e.Logger.Fatal(e.Start(":8080"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
