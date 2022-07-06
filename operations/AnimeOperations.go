package operations

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	"mal/db"
	"mal/models"
	"mal/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type AnimeList struct {
	Data   []models.Data `json:"data"`
	Paging models.Paging `json:"paging"`
}

var (
	MalClientID = models.GetMalClientID()

	GetAnimeQuery = "https://api.myanimelist.net/v2/anime"

	Fields = "fields=id,title,main_picture,alternative_titles,start_date,end_date," +
		"synopsis,mean,rank,popularity,num_list_users,num_scoring_users,nsfw,created_at," +
		"updated_at,media_type,status,genres,my_list_status,num_episodes,start_season," +
		"broadcast,source,average_episode_duration,rating,pictures,background,related_anime," +
		"recommendations,studios,statistics" //TODO: add related_manga field when RelatedManga struct will be added
)

var animeList AnimeList
var client = db.GetRedisClient()

func GetAnimeByID(ID string) models.Node {
	node, err := client.Get(ID).Result()
	handleErr(err)
	var anime models.Node

	err = json.Unmarshal([]byte(node), &anime)
	handleErr(err)
	return anime
}

func GetAnimeRankingList(limit string, offset string) AnimeList {
	req, err := http.NewRequest("GET", GetAnimeQuery+"/ranking"+"?rankingType=all"+"&limit="+limit+"&offset="+offset+"&"+Fields, nil)
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

func GetRandomizedAnime(form *models.RandomizerForm) models.Node {
	animeListLen, err := client.LLen("animeList").Result()
	handleErr(err)
	wholeAnimeList, err := client.LRange("animeList", 0, animeListLen).Result()
	handleErr(err)

	var sortedNodeList []models.Node

	meanFrom, _ := strconv.ParseFloat(form.MeanScoreFrom, 32)
	meanTo, _ := strconv.ParseFloat(form.MeanScoreTo, 32)
	numEpisodesFrom, _ := strconv.Atoi(form.NumEpisodesFrom)
	numEpisodesTo, _ := strconv.Atoi(form.NumEpisodesTo)
	yearFrom, _ := strconv.Atoi(form.YearFrom)
	yearTo, _ := strconv.Atoi(form.YearTo)
	durationFrom, _ := strconv.Atoi(form.DurationFrom)
	durationTo, _ := strconv.Atoi(form.DurationTo)

	for i := 0; int64(i) < animeListLen; i++ {
		var node models.Node
		err = json.Unmarshal([]byte(wholeAnimeList[i]), &node)
		handleErr(err)

		if utils.ContainsOneOfGenres(node.Genres, form.Genres) &&
			(node.Mean >= meanFrom && node.Mean <= meanTo) &&
			(node.NumEpisodes >= numEpisodesFrom && node.NumEpisodes <= numEpisodesTo) &&
			(node.StartSeason.Year >= yearFrom && node.StartSeason.Year <= yearTo) &&
			(node.AverageEpisodeDuration >= durationFrom*60 && node.AverageEpisodeDuration <= durationTo*60) {
			sortedNodeList = append(sortedNodeList, node)
		}
	}

	rand.Seed(time.Now().UnixNano())

	anime := sortedNodeList[rand.Intn(len(sortedNodeList))]

	return anime
}

func GetWholeAnimeList() []AnimeList {
	var list []AnimeList
	list = append(list, GetAnimeRankingList("500", "0"))
	for list[len(list)-1].Paging.Next != "" {
		list = append(list, GetAnimeRankingListViaPaging(list[len(list)-1].Paging))
	}
	return list
}

func GetAnimeListByPage(page string) []models.Node {
	var nodeList []models.Node
	limit, err := strconv.Atoi(page)
	handleErr(err)
	for i := (limit - 1) * 60; i < limit*60; i++ {
		result, err := client.LIndex("animeList", int64(i)).Result()
		handleErr(err)
		var node models.Node
		err = json.Unmarshal([]byte(result), &node)
		handleErr(err)

		nodeList = append(nodeList, node)
	}

	return nodeList
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func DumpAnimeList() {
	list := GetWholeAnimeList()
	var dataList []models.Node
	for _, dataArray := range list {
		for _, data := range dataArray.Data {
			node, err := MarshalBinary(data.Node)
			handleErr(err)
			dataList = append(dataList, data.Node)
			client.Set(strconv.Itoa(data.Node.ID), node, 0)
			client.RPush("animeList", node)
			for _, genre := range data.Node.Genres {
				genreBinary, err := MarshalBinary(genre)
				handleErr(err)
				client.ZAdd("genres", redis.Z{Score: float64(genre.ID), Member: genreBinary})
			}
		}
	}
}

func GetGenres() []models.Genres {
	genresLen, err := client.ZCard("genres").Result()
	handleErr(err)
	result, err := client.ZRange("genres", 0, genresLen).Result()
	handleErr(err)
	var genres []models.Genres
	for _, str := range result {
		var genre models.Genres
		err = json.Unmarshal([]byte(str), &genre)
		handleErr(err)
		genres = append(genres, genre)
	}
	return genres
}

func MarshalBinary(anime interface{}) ([]byte, error) {
	return json.Marshal(anime)
}
