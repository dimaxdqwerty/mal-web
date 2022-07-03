package operations

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mal/models"
	"net/http"
	"strconv"
)

type AnimeList struct {
	Data   []models.Data `json:"data"`
	Paging models.Paging `json:"paging"`
}

var (
	MalClientID = models.GetMalClientID()

	GetAnimeQuery = "https://api.myanimelist.net/v2/anime"

	Fields = "?fields=id,title,main_picture,alternative_titles,start_date,end_date," +
		"synopsis,mean,rank,popularity,num_list_users,num_scoring_users,nsfw,created_at," +
		"updated_at,media_type,status,genres,my_list_status,num_episodes,start_season," +
		"broadcast,source,average_episode_duration,rating,pictures,background,related_anime," +
		"recommendations,studios,statistics" //TODO: add related_manga field when RelatedManga struct will be added
)

var anime models.Anime
var animeList AnimeList

func GetAnimeByID(ID string) models.Anime {
	req, err := http.NewRequest("GET", GetAnimeQuery+"/"+ID+Fields, nil)
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

func GetAnimeRankingList(limit string, offset string) AnimeList {
	req, err := http.NewRequest("GET", GetAnimeQuery+"/ranking"+"?rankingType=all"+"&limit="+limit+"&offset="+offset+Fields, nil)
	req.Header.Add("X-MAL-CLIENT-ID", MalClientID)
	handleErr(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	handleErr(err)

	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)

	err = json.Unmarshal(body, &animeList)
	handleErr(err)
	return animeList
}

func GetAnimeRankingListViaPaging(paging models.Paging) AnimeList {
	var listViaPaging AnimeList
	req, err := http.NewRequest("GET", paging.Next, nil)
	req.Header.Add("X-MAL-CLIENT-ID", MalClientID)
	handleErr(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	handleErr(err)

	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)

	err = json.Unmarshal(body, &listViaPaging)
	handleErr(err)
	return listViaPaging
}

func GetRandomizedAnime() models.Anime {
	return models.Anime{}
}

func GetWholeAnimeList() []AnimeList {
	var list []AnimeList
	list = append(list, GetAnimeRankingList("500", "0"))
	for list[len(list)-1].Paging.Next != "" {
		list = append(list, GetAnimeRankingListViaPaging(list[len(list)-1].Paging))
	}
	return list
}

func GetAnimeListByPage(page string) AnimeList {
	limit, err := strconv.Atoi(page)
	handleErr(err)

	var offset string
	if limit == 1 {
		offset = "0"
	} else {
		offset = strconv.Itoa((limit - 1) * 50)
	}
	return GetAnimeRankingList("50", offset)
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
