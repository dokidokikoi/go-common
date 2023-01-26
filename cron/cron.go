package cron

import (
	"fmt"
	"strconv"
	"time"

	cron3 "github.com/robfig/cron/v3"
)

type TaskFunc func()

type Cron interface {
	AddFunc(fre string, function TaskFunc) (string, error)
	AddDailyFunc(function TaskFunc) (string, error)
	AddHourlyFunc(function TaskFunc) (string, error)
	AddEveryFunc(duration time.Duration, function TaskFunc) (string, error)
	RemoveFunc(id string) error
	StartCron()
	StopCron()
}

type cron struct {
	cr *cron3.Cron
}

func (c *cron) AddFunc(fre string, function TaskFunc) (string, error) {
	identifier, err := c.cr.AddFunc(fre, function)
	return strconv.FormatInt(int64(identifier), 10), err
}

func (c *cron) AddDailyFunc(function TaskFunc) (string, error) {
	identifier, err := c.cr.AddFunc("@daily", function)
	return strconv.FormatInt(int64(identifier), 10), err
}

func (c *cron) AddHourlyFunc(function TaskFunc) (string, error) {
	identifier, err := c.cr.AddFunc("@hourly", function)
	return strconv.FormatInt(int64(identifier), 10), err
}

func (c *cron) AddEveryFunc(duration time.Duration, function TaskFunc) (string, error) {
	spec := fmt.Sprintf("@every %s", duration.String())
	identifier, err := c.cr.AddFunc(spec, function)
	return strconv.FormatInt(int64(identifier), 10), err
}

func (c *cron) RemoveFunc(id string) error {
	identifier, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	c.cr.Remove(cron3.EntryID(identifier))
	return nil
}

func (c *cron) StartCron() {
	c.cr.Start()
}

func (c *cron) StopCron() {
	c.cr.Stop()
}

func NewCronInstance() Cron {
	return &cron{cron3.New()}
}
