package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting GPON Parser - GO Edition")
	runCronJobs()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/gpon/opticalInfo", servOpticalInfo)
	router.GET("/gpon/deviceInfo", servDeviceInfo)
	router.GET("/gpon/allInfo", servAllInfo)

	router.Run("0.0.0.0:8080")
}
