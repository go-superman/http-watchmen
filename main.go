package main

import (
	"github.com/go-superman/http-watchmen/cmd"
	//"github.com/go-superman/http-watchmen/logger"
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
	//logger.Debug("ser done")
}
