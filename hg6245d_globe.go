package main

import (
	"fmt"
	"strconv"

	"github.com/antchfx/htmlquery"
)

func (o HG6245D_Globe) GetGponUrl() string {
	return "http://globebroadband.net/login_new_globe.asp"
}

// cron job
func (o HG6245D_Globe) UpdateCachedPage() {
	doc, err := htmlquery.LoadURL(gponSvc.GetGponUrl())
	if err == nil {
		cachedPage.SetPage(doc)
	} else {
		cachedPage.SetPage(nil)
	}
}

func (o HG6245D_Globe) GetOpticalInfo() *OpticalStats {
	parsedList := make([]string, 0, 5)
	for i := 1; i < 6; i++ {
		htmlNode := htmlquery.FindOne(cachedPage.GetPage(), fmt.Sprintf("/html/body/div[1]/div[1]/div[2]/div/div[4]/ul/li[%d]/span[2]", i))
		parsedList = append(parsedList, normalizeString(htmlquery.InnerText(htmlNode)))
	}

	var opticalInfo = new(OpticalStats)
	opticalInfo.TxPower, _ = strconv.ParseFloat(removeLastNChars(parsedList[0], 4), 64)
	opticalInfo.RxPower, _ = strconv.ParseFloat(removeLastNChars(parsedList[1], 4), 64)
	opticalInfo.Temperature, _ = strconv.ParseFloat(removeLastNChars(parsedList[2], 4), 64)
	opticalInfo.SupplyVoltage, _ = strconv.ParseFloat(removeLastNChars(parsedList[3], 2), 64)
	opticalInfo.BiasCurrent, _ = strconv.ParseFloat(removeLastNChars(parsedList[4], 3), 64)

	return opticalInfo
}

func (o HG6245D_Globe) GetDeviceInfo() *DeviceStats {
	parsedList := make([]string, 0, 6)
	for i := 1; i < 7; i++ {
		htmlNode := htmlquery.FindOne(cachedPage.GetPage(), fmt.Sprintf("/html/body/div[1]/div[1]/div[2]/div/div[5]/ul/li[%d]/span[2]", i))
		parsedList = append(parsedList, normalizeString(htmlquery.InnerText(htmlNode)))
	}

	var deviceInfo = new(DeviceStats)
	deviceInfo.DeviceModel = parsedList[0]
	deviceInfo.ModelSerial = parsedList[1]
	deviceInfo.SoftwareVersion = parsedList[2]
	deviceInfo.MemoryUsage, _ = strconv.ParseFloat(removeLastNChars(parsedList[3], 1), 64)
	deviceInfo.CpuUsage, _ = strconv.ParseFloat(removeLastNChars(parsedList[4], 1), 64)
	deviceInfo.Uptime = parseDuration(parsedList[5])

	return deviceInfo
}
