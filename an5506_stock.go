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
		log.Println(err)
		cachedPage.SetPage(nil)
		cachedPage2.SetPage(nil)
		return
	}

	client := http.Client{Jar: jar}
	form := url.Values{}
	form.Add("User", "user")          // webgui username
	form.Add("Passwd", "tattoo@home") // webgui password

	// 1. Login to UI
	req, err := http.NewRequest("POST", o.GetGponUrl()+"/goform/webLogin", strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(err)
		cachedPage.SetPage(nil)
		cachedPage2.SetPage(nil)
		return
	}

	req.Header.Set("Referer", "http://127.0.0.1/gpon-parser") // must define a referer or it will fail

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		cachedPage.SetPage(nil)
		cachedPage2.SetPage(nil)
		return
	}
	defer resp.Body.Close()

	// 2. Get optical power
	parsePage(client, gponSvc.GetGponUrl()+"/state/opt_power.asp", cachedPage)

	// 3. Get device information
	parsePage(client, gponSvc.GetGponUrl()+"/state/deviceInfor.asp", cachedPage2)

	// 4. Logout from UI
	req, err = http.NewRequest("GET", o.GetGponUrl()+"/goform/webLogout", nil)
	if err != nil {
		log.Println(err)
		cachedPage.SetPage(nil)
		cachedPage2.SetPage(nil)
		return
	}

	req.Header.Set("Referer", "http://192.168.254.254/menu_ph_globe.asp") // must define a referer to logout cleanly
	resp, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
}

func parsePage(client http.Client, url string, docPage *DocPage) {
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		log.Println(err)
	}

	doc, err := htmlquery.Parse(resp.Body)
	if err == nil {
		docPage.SetPage(doc)
	} else {
		log.Println(err)
		docPage.SetPage(nil)
	}

	defer resp.Body.Close()
}

func (o AN5506_Stock) GetOpticalInfo() *OpticalStats {
	var opticalInfo *OpticalStats

	if cachedPage.GetPage() != nil {
		opticalInfo = new(OpticalStats)

		parsedList := make([]string, 0, 5)
		for i := 2; i < 7; i++ {
			htmlNode := htmlquery.FindOne(cachedPage.GetPage(), fmt.Sprintf("/html/body/form/table[3]/tbody/tr[%d]/td[2]", i))
			parsedList = append(parsedList, normalizeString(htmlquery.InnerText(htmlNode)))
		}
		opticalInfo.TxPower, _ = strconv.ParseFloat(removeLastNChars(parsedList[0], 4), 64)
		opticalInfo.RxPower, _ = strconv.ParseFloat(removeLastNChars(parsedList[1], 4), 64)
		opticalInfo.Temperature, _ = strconv.ParseFloat(removeLastNChars(parsedList[2], 4), 64)
		opticalInfo.SupplyVoltage, _ = strconv.ParseFloat(removeLastNChars(parsedList[3], 2), 64)
		opticalInfo.BiasCurrent, _ = strconv.ParseFloat(removeLastNChars(parsedList[4], 3), 64)
	}

	return opticalInfo
}

func (o AN5506_Stock) GetDeviceInfo() *DeviceStats {
	var deviceInfo *DeviceStats

	if cachedPage2.GetPage() != nil {
		deviceInfo = new(DeviceStats)

		parsedList := make([]string, 0, 6)
		for i := 2; i < 16; i++ {
			htmlNode := htmlquery.FindOne(cachedPage2.GetPage(), fmt.Sprintf("/html/body/form/table[3]/tbody/tr[%d]/td[2]", i))
			parsedList = append(parsedList, normalizeString(htmlquery.InnerText(htmlNode)))
		}

		deviceInfo.DeviceModel = parsedList[2]
		deviceInfo.ModelSerial = strings.ReplaceAll(parsedList[4], "FFFFFF", "")
		deviceInfo.SoftwareVersion = parsedList[0]
		deviceInfo.MemoryUsage, _ = strconv.ParseFloat(removeLastNChars(parsedList[10], 1), 64)
		deviceInfo.CpuUsage, _ = strconv.ParseFloat(removeLastNChars(parsedList[9], 1), 64)
		deviceInfo.Uptime = parseDuration(parsedList[13])
	}

	return deviceInfo
}
