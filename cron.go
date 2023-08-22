package main

import (
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/go-co-op/gocron"
)

func runCronJobs() {
	updateCachedPage()

	s := gocron.NewScheduler(time.UTC)
	s.Every(25).Seconds().Do(updateCachedPage)
	s.StartAsync()
}

// cron job
func updateCachedPage() {
	doc, err := htmlquery.LoadURL(GponUrl)
	if err == nil {
		cachedPage.SetPage(doc)
	} else {
		cachedPage.SetPage(nil)
	}
}
