package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Anime struct {
	ID                     int               `json:"id"`
	Title                  string            `json:"title"`
	MainPicture            MainPicture       `json:"main_picture"`
	AlternativeTitles      AlternativeTitles `json:"alternative_titles"`
	StartDate              string            `json:"start_date"`
	EndDate                string            `json:"end_date"`
	Synopsis               string            `json:"synopsis"`
	Mean                   float64           `json:"mean"`
	Rank                   int               `json:"rank"`
	Popularity             int               `json:"popularity"`
	NumListUsers           int               `json:"num_list_users"`
	NumScoringUsers        int               `json:"num_scoring_users"`
	Nsfw                   string            `json:"nsfw"`
	CreatedAt              time.Time         `json:"created_at"`
	UpdatedAt              time.Time         `json:"updated_at"`
	MediaType              string            `json:"media_type"`
	Status                 string            `json:"status"`
	Genres                 []Genres          `json:"genres"`
	MyListStatus           MyListStatus      `json:"my_list_status"`
	NumEpisodes            int               `json:"num_episodes"`
	StartSeason            StartSeason       `json:"start_season"`
	Broadcast              Broadcast         `json:"broadcast"`
	Source                 string            `json:"source"`
	AverageEpisodeDuration int               `json:"average_episode_duration"`
	Rating                 string            `json:"rating"`
	Pictures               []Pictures        `json:"pictures"`
	Background             string            `json:"background"`
	RelatedAnime           []RelatedAnime    `json:"related_anime"`
	//RelatedManga    		[]RelatedManga 		`json:"related_manga"` TODO: create RelatedManga struct
	Recommendations []Recommendations `json:"recommendations"`
	Studios         []Studios         `json:"studios"`
	Statistics      Statistics        `json:"statistics"`
}

var (
	MalClientID = GetMalClientID()

	GetAnimeByIDQuery = "https://api.myanimelist.net/v2/anime/"

	Fields = "?fields=id,title,main_picture,alternative_titles,start_date,end_date," +
		"synopsis,mean,rank,popularity,num_list_users,num_scoring_users,nsfw,created_at," +
		"updated_at,media_type,status,genres,my_list_status,num_episodes,start_season," +
		"broadcast,source,average_episode_duration,rating,pictures,background,related_anime," +
		"recommendations,studios,statistics" //TODO: add related_manga field when RelatedManga struct will be added
)

var anime Anime

func GetAnimeByID(ID string) Anime {
	req, err := http.NewRequest("GET", GetAnimeByIDQuery+ID+Fields, nil)
	req.Header.Add("X-MAL-CLIENT-ID", MalClientID)
	handleErr(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	handleErr(err)

	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)

	err = json.Unmarshal(body, &anime)
	handleErr(err)
	return anime
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
