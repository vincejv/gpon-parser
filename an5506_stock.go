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
	host := getenv("ONT_WEB_HOST", "globebroadband.net")
	port := getenv("ONT_WEB_PORT", "80")
	webProtocol := getenv("ONT_WEB_PROTOCOL", "http")
	return fmt.Sprintf("%s://%s:%s", webProtocol, host, port)
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
	form.Add("User", getenv("ONT_WEB_USER", "user"))          // webgui username
	form.Add("Passwd", getenv("ONT_WEB_PASS", "tattoo@home")) // webgui password

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
	o.parsePage(client, gponSvc.GetGponUrl()+"/state/opt_power.asp", cachedPage)

	// 3. Get device information
	o.parsePage(client, gponSvc.GetGponUrl()+"/state/deviceInfor.asp", cachedPage2)

	// 4. Logout from UI
	req, err = http.NewRequest("GET", o.GetGponUrl()+"/goform/webLogout", nil)
	if err != nil {
		log.Println(err)
		cachedPage.SetPage(nil)
		cachedPage2.SetPage(nil)
		return
	}

	req.Header.Set("Referer", o.GetGponUrl()+"/menu_ph_globe.asp") // must define a referer to logout cleanly
	resp, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
}

func (o AN5506_Stock) parsePage(client http.Client, url string, docPage *DocPage) {
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
			if parseHtmlPage(&parsedList, cachedPage, fmt.Sprintf("/html/body/form/table[3]/tbody/tr[%d]/td[2]", i)) != nil {
				return opticalInfo
			}
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
			if parseHtmlPage(&parsedList, cachedPage2, fmt.Sprintf("/html/body/form/table[3]/tbody/tr[%d]/td[2]", i)) != nil {
				return deviceInfo
			}
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
