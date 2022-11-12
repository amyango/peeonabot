package main

import(
	"time"
	"net/http"
	"log"
	"os"
	"io"
	"encoding/json"
	"strings"
	"strconv"
)

var (
	baseDirectory = "./"
	url string = "https://api.themoviedb.org/3"
	mtokenFile string = baseDirectory + "credentials/moviedb.token"
	token string
)

type MoviePage struct{
	Page int `json:"page"`
	Total_pages int `json:"total_pages"`
	Total_results int `json:"total_results"`	
	Results []Movie `json:"results"`
}

type Movie struct{
	Original_title string `json:"original_title"`
	Release_date string `json:"release_date"`
	Popularity float64 `json:"popularity"`
	Vote_average float64 `json:"vote_average"`
	Vote_count float64 `json:"vote_count"`
	Poster_path string `json:"poster_path"`
}

func (m MoviePage) Print() {
	log.Printf("Page#:         %d ", m.Page)
	log.Printf("Total pages:   %d ", m.Total_pages)
	log.Printf("Total results: %d ", m.Total_results)
}

func (m Movie) Print() {
	log.Printf("[%s] %s (%v)", m.Release_date, m.Original_title, m.Vote_average)
}

func getPage(query string, page int) MoviePage {

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url + "/search/movie", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("api_key", token)
	q.Add("query", query)
	q.Add("page", "" + strconv.Itoa(page))
	req.URL.RawQuery = q.Encode()

	log.Println(q.Encode())

	resp, err := client.Do(req)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)

	var p MoviePage
	json.Unmarshal(body, &p)
	if err != nil {
		log.Panic(err)
	}

	return p
}

// Takes a search string and returns an array of movies
func Find_movies(search []string) []Movie {

	var m []Movie

	if len(search) < 1 {
		return nil
	}

	query := search[0]
	for i := 1; i < len(search); i++ {
		query = query + "+" + search[i]
	}

	page := getPage(query, 1)	
	m = page.Results

	for ;page.Page < page.Total_pages; {
		page.Print()
		log.Printf("before: %d %d\n", page.Page, page.Total_pages)
		page = getPage(query, page.Page + 1)
		log.Printf("after: %d %d\n", page.Page, page.Total_pages)
		m = append(m, page.Results...)
	}

	return m
}

func init() {
	file, err := os.ReadFile(mtokenFile)
	if err != nil {
		log.Panic(err)
	}

	token = strings.Trim(string(file), "\n")
}
