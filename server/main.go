package server

import (
	"github.com/gin-gonic/gin"
	"github.com/dashotv/scry/search"
)

var client *search.Client

func Start(url string) error {
	var err error

	client, err = search.New(url)
	if err != nil {
		return err
	}

	router := gin.Default()
	router.GET("/", homeIndex)
	router.GET("/releases", releasesSearch)

	router.Run(":8080")
	return nil
}
