package core

import (
	"bytes"
	"net/http"

	"github.com/JialeHao/loki-rsyslog-plugin/utils"
)

func workshop(lokicfg *utils.Loki) {
	defer coreWaitGroup.Done()
	for {
		LogData := <-logMsgPool
		go push(LogData, lokicfg)
	}
}

func push(log *LogMsg, lc *utils.Loki) {
	var data []byte
	var httpReq *http.Request
	var httpResp *http.Response
	var err error

	if data, err = packline(log); err != nil {
		logger.Error(err)
		return
	}

	client := &http.Client{}

	if httpReq, err = http.NewRequest("POST", lc.Url, bytes.NewReader(data)); err != nil {
		logger.Error(err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Content-Encoding", "gzip")
	if httpResp, err = client.Do(httpReq); err != nil {
		logger.Error(err)
		return
	}

	respCode := httpResp.StatusCode

	switch respCode / 100 {
	case 2:
		logger.Infof("push success, status_code: %d, proto: %v, ipv4: %v log: %v", respCode, log.proto, log.ip, string(log.msg))
	case 4, 5:
		logger.Errorf("push fatal, error_code=%d", respCode)
	default:
		logger.Warnf("push unknown, status_code: %d", respCode)
	}
}
