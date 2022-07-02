package main

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"log"
	"mal/operations"
	"net/http"
)

var router *gin.Engine

func main() {
	router = gin.Default()
	router.LoadHTMLGlob("src/*")
	router.Use(favicon.New("src/favicon.ico"))

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
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
