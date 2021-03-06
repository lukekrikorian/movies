package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"movies/letterboxd"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	querystring "github.com/google/go-querystring/query"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

type Data struct {
	Count      int16 `json:"movie_count"`
	Limit      int16
	PageNumber int16 `json:"page_number"`
	Movies     []Movie
}

type Response struct {
	Status        string
	StatusMessage string `json:"status_message"`
	Data          Data
}

type Movie struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Title    string
	Rating   float32
	Year     int16
	Runtime  int16
	Summary  string
	Genres   []string
	Language string
	Torrents []Torrent
}

type Torrent struct {
	Hash    string
	Size    string
	Quality string
}

type Query struct {
	Limit          int    `url:"limit,omitempty"`
	Quality        string `url:"quality,omitempty"`
	MinimumRating  int    `url:"minimum_rating,omitempty"`
	QueryTerm      string `url:"query_term,omitempty"`
	Genre          string `url:"genre,omitempty"`
	Sort           string `url:"sort_by,omitempty"`
	Order          string `url:"order_by,omitempty"`
	IncludeRatings bool   `url:"with_rt_ratings,omitempty"`
}

const (
	base = "https://yts.mx/api/v2/list_movies.json"
)

var (
	query           = new(Query)
	disableTrackers = false
	open            = false
	preview         = false
	watchlist       = ""
	trackers        = [...]string{
		"udp://open.demonii.com:1337/announce",
		"udp://tracker.openbittorrent.com:80",
		"udp://tracker.coppersurfer.tk:6969",
		"udp://glotorrents.pw:6969/announce",
		"udp://tracker.opentrackr.org:1337/announce",
		"udp://torrent.gresille.org:80/announce",
		"udp://p4p.arenabg.com:1337",
		"udp://tracker.leechers-paradise.org:6969",
	}
)

var client = http.Client{
	Timeout: time.Second * 10,
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func (q *Query) Search() []Movie {
	queryvalues, err := querystring.Values(q)
	handle(err)
	querystring := queryvalues.Encode()
	url := base + "?" + querystring
	req, err := http.NewRequest(http.MethodGet, url, nil)
	handle(err)

	res, err := client.Do(req)
	handle(err)

	defer res.Body.Close()

	response := new(Response)
	json.NewDecoder(res.Body).Decode(response)
	handle(err)

	return response.Data.Movies
}

func (m *Movie) Magnet() string {
	var magnet string
	for _, torrent := range m.Torrents {
		if torrent.Quality == query.Quality {
			escaped := url.QueryEscape(m.Title)
			spaced := strings.ReplaceAll(escaped, "+", "%20")
			magnet = "magnet:?xt=urn:btih:" + torrent.Hash + "&dn=" + spaced
		}
	}
	if !disableTrackers && magnet != "" {
		for _, tracker := range trackers {
			magnet += "&tr=" + tracker
		}
	}
	return magnet
}

func main() {
	app := &cli.App{
		Name:  "movies",
		Usage: "search yts.mx for torrents",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "query",
				Aliases:     []string{"q"},
				Usage:       "`QUERY` to search",
				Destination: &query.QueryTerm,
			},
			&cli.IntFlag{
				Name:        "rating",
				Aliases:     []string{"r"},
				Usage:       "minimum imdb user `RATING` to filter by: 0 to 9 inclusive",
				Destination: &query.MinimumRating,
			},
			&cli.StringFlag{
				Name:        "quality",
				Aliases:     []string{"qual"},
				Value:       "1080p",
				Usage:       "file `QUALITY` to filter by: 720p, 1080p, 2160p, or 3D",
				Destination: &query.Quality,
			},
			&cli.StringFlag{
				Name:        "genre",
				Aliases:     []string{"g"},
				Usage:       "imdb `GENRE` from https://www.imdb.com/genre/ to filter by",
				Destination: &query.Genre,
			},
			&cli.StringFlag{
				Name:        "sort",
				Aliases:     []string{"s"},
				Usage:       "`VALUE` to sort by: title, year, rating, peers, seeds, download_count, like_count, or date_added",
				Destination: &query.Sort,
			},
			&cli.StringFlag{
				Name:        "order",
				Aliases:     []string{"o"},
				Usage:       "`ORDER` to order results by: desc or asc",
				Destination: &query.Order,
			},
			&cli.BoolFlag{
				Name:        "disable-trackers",
				Aliases:     []string{"dt"},
				Usage:       "disables trackers in generated magnet links",
				Destination: &disableTrackers,
			},
			&cli.BoolFlag{
				Name:        "open",
				Usage:       "opens the first search result magnet link",
				Destination: &open,
			},
			&cli.StringFlag{
				Name:        "watchlist",
				Aliases:     []string{"l"},
				Usage:       "retrieves a random film from a letterboxd watchlist and searches it",
				Destination: &watchlist,
			},
			&cli.BoolFlag{
				Name:        "preview",
				Usage:       "Opens the movie on letterboxd if searching from the watchlist",
				Destination: &preview,
			},
		},
		Action: func(c *cli.Context) error {
			if watchlist != "" {
				movies := letterboxd.Watchlist(watchlist)
				if len(movies) < 1 {
					fmt.Println("No movies in watchlist.")
					return nil
				}

				rand.Seed(time.Now().Unix())
				index := rand.Intn(len(movies))
				movie := movies[index]
				movietitle := fmt.Sprintf("%s %s",
					movie[letterboxd.Title],
					movie[letterboxd.Year])
				query.QueryTerm = movietitle
				fmt.Println("Searching for", movietitle)

				if preview {
					browser.OpenURL(movie[letterboxd.URL])
				}
			}

			movies := query.Search()
			for i, movie := range movies {
				magnet := movie.Magnet()
				if magnet == "" {
					continue
				}

				if i > 0 {
					fmt.Println()
				}

				fmt.Printf("%s (%d)\n-- %s\n", movie.Title, movie.Year, magnet)
				if open && i == 0 {
					browser.OpenURL(magnet)
				}
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
