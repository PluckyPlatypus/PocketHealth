package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"pocket-health/pkg/controller"
)

func main() {
	controller := controller.NewController()
	router := gin.Default()

	// add swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/dicoms", controller.PostDicom)
	router.GET("/dicoms/:studyId", controller.GetDicomPng)
	router.GET("/dicoms/:studyId/:tag", controller.GetDicomAttributeHeader)

	err := router.Run("localhost:8080")

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}
