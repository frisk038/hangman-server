package cron

import "github.com/robfig/cron/v3"

func NewCronMidnight(task func()) *cron.Cron {
	c := cron.New()
	c.AddFunc("@every 0h05m00s", task)
	return c
}
