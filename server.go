package main

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"log"
	"mal/models"
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
		context.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "Home",
			})
	})
	router.GET("/anime/:id", func(context *gin.Context) {
		id := context.Param("id")
		anime := models.GetAnimeByID(id)
		context.HTML(
			http.StatusOK,
			"anime.html",
			gin.H{
				"anime": anime,
			})
	})
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
