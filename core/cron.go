package core

import (
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/vincejv/gpon-parser/device"
	"github.com/vincejv/gpon-parser/util"
)

func RunCronJobs() {
	s := gocron.NewScheduler(time.UTC)
	pollTime, _ := strconv.Atoi(util.Getenv("ONT_POLL_SEC", "60")) // ignore error, default to 25 on failure
	s.Every(pollTime).Seconds().Do(device.GponSvc.UpdateCachedPage)
	s.StartAsync()
}
