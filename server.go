package main

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"log"
	"mal/operations"
	"net/http"
	"strconv"
)

var router *gin.Engine

func main() {
	router = gin.Default()
	router.LoadHTMLGlob("src/*")
	router.Use(favicon.New("src/favicon.ico"))
	router.Static("/assets", "./assets")

	initializeRoutes()

	err := router.Run("localhost:8080")
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
		context.HTML(
			http.StatusOK,
			"anime.html",
			gin.H{
				"anime": anime,
			})
	})
	/*router.POST("randomize", func(context *gin.Context) {
		err := context.Request.ParseForm()
		handleErr(err)
		genres := context.PostFormArray("genres")
		numEpisodes := context.PostForm("numEpisodes")
		fmt.Println(genres)
		fmt.Println(numEpisodes)
		operations.GetRandomizedAnime()
	})*/
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
