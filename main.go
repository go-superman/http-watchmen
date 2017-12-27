package main

import (
	"github.com/hyhlinux/http-watchmen/cmd"
	"github.com/hyhlinux/http-watchmen/logger"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		cmd.Execute()
		wg.Done()
	}()

	wg.Wait()
	logger.Debug("ser done")
}
