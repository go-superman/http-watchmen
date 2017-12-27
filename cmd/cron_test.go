package cmd

import (
	"github.com/go-superman/http-watchmen/logger"
	"github.com/robfig/cron"
	"sync"
	"testing"
	"time"
)

const ONE_SECOND = 1*time.Second + 10*time.Millisecond

func TestCronCall(t *testing.T) {
	logger.Debugf("cronCall..")
	wg := &sync.WaitGroup{}
	c := cron.New()
	// cron.NewWithLocation 使用自定义时区
	//c.AddFunc("*/20 * * * * *", func() {
	c.AddFunc("@every 21s", func() {
		logger.Debug("Every 20s do job")
		wg.Done()
	})
	wg.Add(1)
	c.AddFunc("@every 5s", func() {
		logger.Debug("Every 5s do job")
		//wg.Done()
	})
	//c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
	c.Start()
	// Funcs are invoked in their own goroutine, asynchronously.
	//...
	// Funcs may also be added to a running Cron
	//c.AddFunc("@daily", func() { fmt.Println("Every day") })

	defer c.Stop() // Stop the scheduler (does not stop any jobs already running).
	for {
		select {
		case <-time.After(ONE_SECOND * 5):
			logger.Debugf("display jobs:")
			displayJob(c)
		case <-wait(wg):
			logger.Debugf("job1 done")
			displayJob(c)
			goto FOREND
		}
	}

FOREND:
	logger.Debugf("cronjob end")
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

func stop(cron *cron.Cron) chan bool {
	ch := make(chan bool)
	go func() {
		cron.Stop()
		ch <- true
	}()
	return ch
}
