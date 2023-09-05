package main

import (
	"time"

	"github.com/go-co-op/gocron"
)

func runCronJobs() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(25).Seconds().Do(gponSvc.UpdateCachedPage)
	s.StartAsync()
}
