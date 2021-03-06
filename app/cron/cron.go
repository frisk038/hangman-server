package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

func NewCronMidnight(task func()) *cron.Cron {
	c := cron.New()
	_, err := c.AddFunc("@midnight", task)
	if err != nil {
		log.Fatalf("adding cron task fail : %s", err)
	}
	return c
}
