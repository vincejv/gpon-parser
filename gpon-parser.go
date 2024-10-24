package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

/**

Available ENVIRONMENT Variables

ONT_WEB_HOST
ONT_WEB_PORT
ONT_WEB_PROTOCOL
ONT_WEB_USER
ONT_WEB_PASS
ONT_TELNET_PORT
LISTEN_PORT
LISTEN_IP

**/

func main() {
	log.Println("Starting GPON Parser")

	initGponSvc()
	runCronJobs()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/gpon/opticalInfo", servOpticalInfo)
	router.GET("/gpon/deviceInfo", servDeviceInfo)
	router.GET("/gpon/allInfo", servAllInfo)
	router.GET("/health", health)

	ip := getenv("LISTEN_IP", "0.0.0.0")
	port := getenv("LISTEN_PORT", "8092")
	log.Printf("Starting web server on %s:%s\n", ip, port)
	log.Printf("Polling ONT every %ss", getenv("ONT_POLL_SEC", "60"))
	router.Run(fmt.Sprintf("%s:%s", ip, port))
}

func initGponSvc() {
	telnetInit.SetFlag(false)
	telnetScripts.SetFlag(false)

	model := os.Getenv("ONT_MODEL")
	if len(model) > 0 {
		if strings.EqualFold(model, "an5506_stock") {
			log.Println("ONT Model is Fiberhome AN5506")
			gponSvc = new(AN5506_Stock)
		} else if strings.EqualFold(model, "hg6245d_globe") {
			log.Println("ONT Model is Fiberhome HG6245D")
			gponSvc = new(HG6245D_Globe)
		} else if strings.EqualFold(model, "zte_f670") {
			log.Println("ONT Model is ZTE F670L")
			gponSvc = new(ZTEF670L)
		} else if strings.EqualFold(model, "zlt_g3000a") {
			log.Println("ONT Model is ZLT G3000A WiFi 6")
			gponSvc = new(ZLTG3000A)
		} else {
			log.Println("Invalid ONT model provided in env variable 'ONT_MODEL', valid args are ['an5506_stock', 'hg6245d_globe', 'zte_f670', 'zlt_g3000a']")
			os.Exit(-10)
		}
	} else {
		// by default use AN5506 stock
		log.Println("Did not specify any ONT models in env variable 'ONT_MODEL', using default model Fiberhome AN5506")
		gponSvc = new(AN5506_Stock)
	}
}
