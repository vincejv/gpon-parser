package core

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/vincejv/gpon-parser/device"
	"github.com/vincejv/gpon-parser/util"
)

func RunCronJobs() {
	s := gocron.NewScheduler(time.UTC)
	pollTime := util.ParseInt(util.Getenv("ONT_POLL_SEC", "60")) // ignore error, default to 25 on failure
	s.Every(pollTime).Seconds().Do(device.GponSvc.UpdateCachedPage)
	gocron.SetPanicHandler(func(jobName string, _ interface{}) {
		log.Printf("Panic in job: %s\n", jobName)
	})
	s.StartAsync()
}
