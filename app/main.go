package main

import (
	"log"
	"os"

	v1 "github.com/frisk038/hangman-server/app/handler/v1"
	"github.com/frisk038/hangman-server/infra/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	_, err := repository.NewClient()
	if err != nil {
		log.Fatalf("db init fail %s", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/secret", v1.GetSecret)

	router.Run(":" + port)
}
