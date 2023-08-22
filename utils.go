package main

import (
	"strconv"
	"strings"
)

func normalizeString(s string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

func removeLastNChars(s string, lengthNChars int) string {
	return s[:len(s)-lengthNChars]
}

func parseDuration(timeString string) int64 {
	durationSplit := strings.Fields(timeString)
	daysConv, _ := strconv.ParseInt(durationSplit[0], 10, 64)
	hoursConv, _ := strconv.ParseInt(durationSplit[2], 10, 64)
	minsConv, _ := strconv.ParseInt(durationSplit[4], 10, 64)
	secsConv, _ := strconv.ParseInt(durationSplit[6], 10, 64)
	return daysConv*86400 + hoursConv*3600 + minsConv*60 + secsConv
}
