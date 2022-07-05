package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	"github.com/thinkerou/favicon"
	"log"
	"mal/db"
	"mal/models"
	"mal/operations"
	"net/http"
	"strconv"
)

var router *gin.Engine
var client = db.GetRedisClient()

func main() {
	router = gin.Default()
	router.LoadHTMLGlob("src/*")
	router.Use(favicon.New("src/favicon.ico"))
	router.Static("/assets", "./assets")

	initializeRoutes()
	dumpAnimeListJob()

	go func() {
		err := gocron.Every(2).Hours().Do(dumpAnimeListJob)
		handleErr(err)
		<-gocron.Start()
	}()

	err := router.Run()
	handleErr(err)
}

func initializeRoutes() {
	router.GET("/", func(context *gin.Context) {
		//animeList := operations.GetAnimeRankingList("100", "0")
		context.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				//"list": animeList,
			})
	})
	router.GET("/anime/:id", func(context *gin.Context) {
		id := context.Param("id")
		anime := operations.GetAnimeByID(id)

		averageEpisodeDurationMin := anime.AverageEpisodeDuration / 60
		context.HTML(
			http.StatusOK,
			"anime.html",
			gin.H{
				"anime":                     anime,
				"averageEpisodeDurationMin": averageEpisodeDurationMin,
			})
	})
	router.POST("randomize", func(context *gin.Context) {
		err := context.Request.ParseForm()
		handleErr(err)

		genres := context.PostFormArray("genres")
		meanScoreFrom := context.PostForm("meanScoreFrom")
		meanScoreTo := context.PostForm("meanScoreTo")
		numEpisodesFrom := context.PostForm("numEpisodesFrom")
		numEpisodesTo := context.PostForm("numEpisodesTo")
		yearFrom := context.PostForm("yearFrom")
		yearTo := context.PostForm("yearTo")
		durationFrom := context.PostForm("durationFrom")
		durationTo := context.PostForm("durationTo")

		randomAnime := operations.GetRandomizedAnime(&models.RandomizerForm{
			Genres:          genres,
			MeanScoreFrom:   meanScoreFrom,
			MeanScoreTo:     meanScoreTo,
			NumEpisodesFrom: numEpisodesFrom,
			NumEpisodesTo:   numEpisodesTo,
			YearFrom:        yearFrom,
			YearTo:          yearTo,
			DurationFrom:    durationFrom,
			DurationTo:      durationTo,
		})

		context.Redirect(http.StatusSeeOther, "/anime/"+strconv.Itoa(randomAnime.Node.ID))
	})
	router.GET("/anime/all", func(context *gin.Context) {
		list := operations.GetWholeAnimeList()
		context.HTML(
			http.StatusOK,
			"animeAll.html",
			gin.H{
				"list": list,
			})
	})
	router.GET("/animeList", func(context *gin.Context) {
		page := context.Query("page")
		nextPage, err := strconv.Atoi(page)
		handleErr(err)

		previousPage := nextPage - 1
		if nextPage <= 1 {
			previousPage = 1
		}
		list := operations.GetAnimeListByPage(page)
		context.HTML(
			http.StatusOK,
			"animeList.html",
			gin.H{
				"list":         list,
				"previousPage": previousPage,
				"next1":        nextPage + 1,
				"next2":        nextPage + 2,
				"next3":        nextPage + 3,
				"nextPage":     nextPage + 1,
			})
	})
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func dumpAnimeListJob() {
	operations.DumpAnimeList()
	_, time := gocron.NextRun()
	fmt.Println("Next dump in ", time)
}
