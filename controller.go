package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func servOpticalInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gponSvc.GetOpticalInfo())
}

func servDeviceInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gponSvc.GetDeviceInfo())
}

func servAllInfo(c *gin.Context) {
	allInfo := new(AllStats)

	allInfo.OpticalStats = gponSvc.GetOpticalInfo()
	allInfo.DeviceStats = gponSvc.GetDeviceInfo()

	c.JSON(http.StatusOK, allInfo)
}
