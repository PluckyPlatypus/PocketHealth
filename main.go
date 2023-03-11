package main

import (
	"errors"
	"fmt"
	"net/http"

	"pocket-health/pkg/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	controller := controller.NewController()
	router := gin.Default()

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
