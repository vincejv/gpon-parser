package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting GPON Parser - GO Edition")
	initGponSvc()
	runCronJobs()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/gpon/opticalInfo", servOpticalInfo)
	router.GET("/gpon/deviceInfo", servDeviceInfo)
	router.GET("/gpon/allInfo", servAllInfo)

	router.Run("0.0.0.0:8092")
}

func initGponSvc() {
	if len(os.Args) >= 2 {
		if strings.EqualFold(os.Args[1], "an5506_stock") {
			gponSvc = new(AN5506_Stock)
		} else if strings.EqualFold(os.Args[1], "hg6245d_globe") {
			gponSvc = new(HG6245D_Globe)
		} else {
			log.Println("Invalid ONT model provided in args, valid args are ['an5506_stock', 'hg6245d_globe']")
			os.Exit(-10)
		}
	} else {
		// by default use AN5506 stock
		log.Println("Did not specify any ONT models in CLI args, using default model Fiberhome AN5506")
		gponSvc = new(AN5506_Stock)
	}
}
