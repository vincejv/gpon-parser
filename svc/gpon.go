package svc

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vincejv/gpon-parser/device"
	"github.com/vincejv/gpon-parser/model"
)

func ServOpticalInfo(c *gin.Context) {
	stat := device.GponSvc.GetOpticalInfo()
	if stat != nil && stat.Temperature > 0.0 && stat.RxPower < -5.0 && stat.TxPower > 0.1 {
		c.JSON(http.StatusOK, stat)
	} else {
		log.Println("Unable to fetch gpon optical stats at the moment, returning http 500")
		log.Printf("servOpticalInfo: %+v\n", stat)
		c.JSON(http.StatusInternalServerError, nil)
	}
}

func ServDeviceInfo(c *gin.Context) {
	stat := device.GponSvc.GetDeviceInfo()
	if stat != nil && len(strings.TrimSpace(stat.DeviceModel)) >= 0 {
		c.JSON(http.StatusOK, device.GponSvc.GetDeviceInfo())
	} else {
		log.Println("Unable to fetch gpon device stats at the moment, returning http 500")
		log.Printf("servDeviceInfo: %+v\n", stat)
		c.JSON(http.StatusInternalServerError, nil)
	}
}

func ServAllInfo(c *gin.Context) {
	allInfo := new(model.AllStats)

	allInfo.OpticalStats = device.GponSvc.GetOpticalInfo()
	allInfo.DeviceStats = device.GponSvc.GetDeviceInfo()

	if allInfo.OpticalStats != nil && allInfo.DeviceStats != nil &&
		allInfo.OpticalStats.Temperature > 0.0 && len(strings.TrimSpace(allInfo.DeviceStats.DeviceModel)) >= 0 {
		c.JSON(http.StatusOK, allInfo)
	} else {
		log.Println("Unable to fetch all gpon stats at the moment, returning http 500")
		log.Printf("servAllInfo: %+v\n", allInfo)
		c.JSON(http.StatusInternalServerError, nil)
	}
}
