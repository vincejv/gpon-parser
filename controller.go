package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func servOpticalInfo(c *gin.Context) {
	c.JSON(http.StatusOK, getOpticalInfo())
}

func servDeviceInfo(c *gin.Context) {
	c.JSON(http.StatusOK, getDeviceInfo())
}

func servAllInfo(c *gin.Context) {
	allInfo := new(AllStats)

	allInfo.OpticalStats = getOpticalInfo()
	allInfo.DeviceStats = getDeviceInfo()

	c.JSON(http.StatusOK, allInfo)
}
