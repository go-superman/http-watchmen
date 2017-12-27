package utils

import (
	"fmt"
	"github.com/go-superman/http-watchmen/logger"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

func HealthCheck(url string, retryCnt int, retryTime time.Duration) (err error) {
	// 超过一定时间，返回非200，健康检查失败.
	resp, _, errs := gorequest.New().Get(url).
		Retry(retryCnt, retryTime, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError).
		End()
	defer func() {
		if err != nil {
			retryCountReturn := resp.Header.Get("Retry-Count")
			logger.Warnf("Expected [%v] retry but was [%v]", retryCnt, retryCountReturn)
		}
	}()
	if errs != nil {
		err = fmt.Errorf("HealthCheck: %v resp:%v", errs, resp)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HealthCheck status:%v", resp.StatusCode)
		return err
	}
	logger.Debugf("HealthCheck OK:%v url:%v", resp.StatusCode, url)
	return nil
}
