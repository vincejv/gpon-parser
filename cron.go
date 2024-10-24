package main

import (
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
)

func runCronJobs() {
	s := gocron.NewScheduler(time.UTC)
	pollTime, _ := strconv.Atoi(getenv("ONT_POLL_SEC", "60")) // ignore error, default to 25 on failure
	s.Every(pollTime).Seconds().Do(gponSvc.UpdateCachedPage)
	s.StartAsync()
}
