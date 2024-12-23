package util

import (
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

func NormalizeString(s string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

func RemoveLastNChars(s string, lengthNChars int) string {
	return s[:len(s)-lengthNChars]
}

func ParseDuration(timeString string) int64 {
	durationSplit := strings.Fields(timeString)
	daysConv, _ := strconv.ParseInt(durationSplit[0], 10, 64)
	hoursConv, _ := strconv.ParseInt(durationSplit[2], 10, 64)
	minsConv, _ := strconv.ParseInt(durationSplit[4], 10, 64)
	secsConv, _ := strconv.ParseInt(durationSplit[6], 10, 64)
	return daysConv*86400 + hoursConv*3600 + minsConv*60 + secsConv
}

func RandInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func ParseHtmlPage(elemList *[]string, page *DocPage, xpath string) error {
	htmlNode := htmlquery.FindOne(page.GetPage(), xpath)
	if htmlNode != nil {
		*elemList = append(*elemList, NormalizeString(htmlquery.InnerText(htmlNode)))
		return nil
	}
	return errors.New("unable to find xpath: " + xpath)
}
