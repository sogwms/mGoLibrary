package job

import (
	"context"

	"github.com/robfig/cron/v3"
)

type jobHandle interface{}

var gCron *cron.Cron

func init() {
	gCron = cron.New()
	gCron.Start()
}

func NewJob(schedule string, f func()) (jobHandle, error) {
	return gCron.AddFunc(schedule, f)
}

func RemoveJob(job jobHandle) bool {
	entry, ok := job.(cron.EntryID)
	if ok {
		gCron.Remove(entry)
	}

	return ok
}

func StopAll() context.Context {
	return gCron.Stop()
}
