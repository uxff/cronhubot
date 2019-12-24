package runner

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/rubyist/circuitbreaker"
	"github.com/uxff/cronhubot/pkg/log"
)

type Config struct {
	Url     string
	Retries uint
	Timeout uint
}

// 定时任务执行时，业务方需要返回的应答消息格式
type RemoteResponse struct {
	// 状态码
	// 200 正常，继续执行
	// 400 定时任务已不需要继续执行，可删除
	Code int `json:"code"`

	// 业务方自定义回复消息
	Message string `json:"Message"`
}

// 远端已经不需要本服务继续执行定时任务
func (r *RemoteResponse) UselessJob() bool {
	return r.Code == 400
}

// 执行定时任务，调用远端的接口
func NoticeRemote(traceId string, c *Config) (*RemoteResponse, error) {
	timeout := time.Second * time.Duration(c.Timeout)
	client := circuit.NewHTTPClient(timeout, int64(c.Retries), nil)

	resp, err := client.Get(c.Url)
	if err != nil {
		log.Trace(traceId).Errorf("Failed to send event:%v", err)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Trace(traceId).Errorf("读取应答失败:%v", err)
		return nil, err
	}

	var res = new(RemoteResponse)
	if err := json.Unmarshal(b, res); err != nil {
		log.Trace(traceId).Errorf("解析应答失败:%v 回复：%s", err, b)
		return nil, err
	}

	log.Trace(traceId).Infof("定时任务执行成功, 配置:%+v 回复:%+v", c, res)
	return res, nil
}
