package main

import (
	"log"
	"os"

	v1 "github.com/frisk038/hangman-server/app/handler/v1"
	"github.com/frisk038/hangman-server/business/usecase"
	"github.com/frisk038/hangman-server/infra/adapter/giphy"
	"github.com/frisk038/hangman-server/infra/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Create DB
	repo, err := repository.NewClient()
	if err != nil {
		log.Fatalf("db init fail %s", err)
	}

	// Giphy adapter
	giphy := giphy.NewGiphy()

	// Create business
	ps := usecase.NewProcessSecret(repo)
	gf := usecase.NewProcessGIF(giphy)

	// Create handler
	handlers := v1.NewSecretHandler(ps, gf)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/secret", handlers.GetSecret)
	router.POST("/score", handlers.PostScore)
	router.POST("/user", handlers.UpdateUserName)
	router.GET("/top", handlers.SelectTopUser)
	router.GET("/wingif", handlers.GetSuccessGif)
	router.Run(":" + port)
}
