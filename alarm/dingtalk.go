package alarm

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/JialeHao/loki-rsyslog-plugin/utils"
)

var logger = utils.InitLogger()

var (
	dt   *dingtalk
	once sync.Once
)

type AlarmChannel interface {
	Push(log *[]byte)
}

type dingtalk struct {
	dt utils.Dingtalk
}

type alarmMsg struct {
	Msgtype string `json:"msgtype"`
	Markdown markdown `json:"markdown"`
}

type markdown struct {
	Title string `json:"title"`
	Text string `json:"text"`
}

func (dingtalk *dingtalk) Push(log *[]byte) {
	var httpReq *http.Request
	var httpResp *http.Response
	var err error
	var msgjson []byte

	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signString := ts + "\n" + dingtalk.dt.Secret
	hmac := hmac.New(sha256.New, []byte(dingtalk.dt.Secret))
	hmac.Write([]byte(signString))
	hmacDigest := hmac.Sum(nil)
	sign := base64.StdEncoding.EncodeToString(hmacDigest)

	url := "https://oapi.dingtalk.com/robot/send?access_token=" + dingtalk.dt.Token + "&timestamp=" + ts + "&sign=" + sign

	client := &http.Client{}

	msg := alarmMsg{
		Msgtype: "markdown",
		Markdown: markdown{
			Title: "日志告警",
			Text: string(*log),
		},
	}

	if msgjson , err = json.Marshal(msg); err != nil {
		logger.Error(err)
		return
	}

	fmt.Println(string(msgjson))

	if httpReq, err = http.NewRequest("POST", url, bytes.NewReader(msgjson)); err != nil {
		logger.Error(err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if httpResp, err = client.Do(httpReq); err != nil {
		logger.Error(err)
		return
	}

	var resbyte []byte

	if resbyte, err = io.ReadAll(httpResp.Body); err != nil {
		logger.Errorf("err_code: %v, return_body: %v", httpResp.StatusCode, string(resbyte))
		return
	}
	logger.Infof("alarm_code: %v, return_body: %v", httpResp.StatusCode, string(resbyte))
}

func newAlarmTool() *dingtalk {
	return &dingtalk{
		dt: utils.GlobalConfig.Dingtalk,
	}
}

func InitAlarmChannel() AlarmChannel {
	once.Do(func() {
		dt = newAlarmTool()
	})

	return dt
}
