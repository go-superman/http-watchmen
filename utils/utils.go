package utils

import (
	"fmt"
	"encoding/json"
	"github.com/go-superman/http-watchmen/logger"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)



func HealthCheck(url string, retryCnt int, status []int, timeout, retryTime time.Duration) (data string, err error) {
	// 超过一定时间，返回非200，健康检查失败.
	dataInfo := make(map[string]interface{})
	dataInfo["url"] = url
	dataInfo["create_time"] = time.Now().UTC().String()
	if timeout <= 0 {
		timeout = 5*time.Second
	}
	if len(status) == 0 {
		status = []int{http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError}
	}
	resp, _, errs := gorequest.New().Timeout(timeout).Get(url).
		Retry(retryCnt, retryTime, status...).
		End()
	for i:=0; i<retryCnt; i++ {
		resp, _, errs := gorequest.New().Get(url).End()
		logger.Warnf("resp:%v errs:%v", resp.StatusCode, errs)
	}
	defer func() {
		dataInfo["retry_count_return"] = 0
		dataInfo["err"] = ""
		dataInfo["end_time"] = time.Now().UTC().String()
		if err != nil {
			retryCountReturn := ""
			if resp != nil {
				retryCountReturn = resp.Header.Get("Retry-Count")
				dataInfo["retry_count_return"] = retryCountReturn
			}
			dataInfo["err"] = fmt.Sprintf("%v", err)
			logger.Warnf("Expected [%v] retry but was [%v] url:%v", retryCnt, retryCountReturn, url)
		}
		if dataByte, e := json.Marshal(dataInfo); e != nil {
			logger.Errorf("e:%v datainfo:%v", e, dataInfo)
		}else{
			data = string(dataByte)
		}

	}()
	if errs != nil {
		err = fmt.Errorf("HealthCheck: %v resp:%v url:%v", errs, resp, url)
		return data, err
	}

	dataInfo["status_code"] = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HealthCheck status:%v url:%v", resp.StatusCode, url)
		return data, err
	}
	logger.Debugf("HealthCheck OK:%v url:%v", resp.StatusCode, url)
	return data,nil
}
