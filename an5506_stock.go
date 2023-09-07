package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/publicsuffix"
)

func (o AN5506_Stock) GetGponUrl() string {
	return "http://globebroadband.net"
}

// cron job
func (o AN5506_Stock) UpdateCachedPage() {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}

	form := url.Values{}
	form.Add("User", "user")          // webgui username
	form.Add("Passwd", "tattoo@home") // webgui password
	req, _ := http.NewRequest("POST", o.GetGponUrl()+"/goform/webLogin", strings.NewReader(form.Encode()))
	req.Header.Set("Referer", "http://127.0.0.1/gpon-parser") // must define a referer or it will fail
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	parsePage(client, gponSvc.GetGponUrl()+"/state/opt_power.asp", cachedPage)
	parsePage(client, gponSvc.GetGponUrl()+"/state/deviceInfor.asp", cachedPage2)

	req, _ = http.NewRequest("GET", o.GetGponUrl()+"/goform/webLogout", nil)
	req.Header.Set("Referer", "http://192.168.254.254/menu_ph_globe.asp") // must define a referer to logout cleanly
	client.Do(req)
}

func parsePage(client http.Client, url string, docPage *DocPage) {
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	doc, err := htmlquery.Parse(resp.Body)
	if err == nil {
		docPage.SetPage(doc)
	} else {
		log.Fatal(err)
		docPage.SetPage(nil)
	}
}

func (o AN5506_Stock) GetOpticalInfo() *OpticalStats {
	parsedList := make([]string, 0, 5)
	for i := 2; i < 7; i++ {
		htmlNode := htmlquery.FindOne(cachedPage.GetPage(), fmt.Sprintf("/html/body/form/table[3]/tbody/tr[%d]/td[2]", i))
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

func (o AN5506_Stock) GetDeviceInfo() *DeviceStats {
	parsedList := make([]string, 0, 6)
	for i := 2; i < 16; i++ {
		htmlNode := htmlquery.FindOne(cachedPage2.GetPage(), fmt.Sprintf("/html/body/form/table[3]/tbody/tr[%d]/td[2]", i))
		parsedList = append(parsedList, normalizeString(htmlquery.InnerText(htmlNode)))
	}

	var deviceInfo = new(DeviceStats)
	deviceInfo.DeviceModel = parsedList[2]
	deviceInfo.ModelSerial = strings.ReplaceAll(parsedList[4], "FFFFFF", "")
	deviceInfo.SoftwareVersion = parsedList[0]
	deviceInfo.MemoryUsage, _ = strconv.ParseFloat(removeLastNChars(parsedList[10], 1), 64)
	deviceInfo.CpuUsage, _ = strconv.ParseFloat(removeLastNChars(parsedList[9], 1), 64)
	deviceInfo.Uptime = parseDuration(parsedList[13])

	return deviceInfo
}
