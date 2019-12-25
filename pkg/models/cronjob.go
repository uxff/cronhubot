package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron"
	"github.com/uxff/cronhubot/pkg/utils"
)

const (
	StatusActive   = 1
	StatusInactive = 99
	DefaultTimeStr = "1999-01-01 00:00:00"
)

type CronJob struct {
	Id             uint           `json:"id"               xorm:"unsigned int(11) notnull pk autoincr id"`
	Name           string         `json:"name"             xorm:"varchar(50) default('') 'name'"`
	Url            string         `json:"url"              xorm:"varchar(256) notnull default('') 'url'"`
	Expression     string         `json:"expression"       xorm:"varchar(100) notnull default('') 'expression'"`
	Status         uint           `json:"status"           xorm:"tinyint(10) default(1) 'status'"`
	Retries        uint           `json:"retries"          xorm:"unsigned int(11) default(0) 'retries'"`
	RequestTimeout uint           `json:"request_timeout"  xorm:"unsigned int(11) default(3) 'request_timeout'"`
	StopAt         utils.JsonTime `json:"stop_at"          xorm:"timestamp not null default('1999-01-01 00:00:00') 'stop_at'"` // 最终停止时间
	CreatedAt      utils.JsonTime `json:"created_at"       xorm:"timestamp not null default('1999-01-01 00:00:00') 'created_at'"`
	UpdatedAt      utils.JsonTime `json:"updated_at"       xorm:"timestamp not null default('1999-01-01 00:00:00') 'updated_at'"`
}

func NewCronJob() *CronJob {
	return &CronJob{
		Status:         StatusActive, // 默认是有效的定时任务
		Retries:        0,            // 默认请求业务方时不重试
		RequestTimeout: 3,            // 请求业务方时默认超时时间是3s
	}
}

func (e *CronJob) TableName() string {
	return "athena_cronjobs"
}

func (e *CronJob) CheckExpression(expression string) error {
	if len(strings.Split(expression, " ")) != 6 {
		return errors.New("定时配置格式不正确")
	}

	if _, err := cron.Parse(expression); err != nil {
		return err
	}
	return nil
}

func (e *CronJob) Validate() (errors map[string]string, ok bool) {
	errors = make(map[string]string)
	if e.Url == "" {
		errors["url"] = "url is empty"
	}

	if err := e.CheckExpression(e.Expression); err != nil {
		errors["expression"] = err.Error()
	}

	if e.Status != StatusActive && e.Status != StatusInactive {
		errors["status"] = fmt.Sprintf("status must be %v or %v", StatusActive, StatusInactive)
	}

	if e.Retries < 0 || e.Retries > 10 {
		errors["retries"] = "field retries must be between 0 and 10"
	}

	if !e.StopAt.IsEmptyTime() {
		if e.StopAt.ToStdTime().Before(time.Now()) {
			errors["stop_at"] = "stop_at before now"
		}
	}

	ok = len(errors) == 0

	return
}

func (e *CronJob) SetAttributes(newEnt *CronJob) {
	if newEnt.Name != "" {
		e.Name = newEnt.Name
	}

	if newEnt.Url != "" {
		e.Url = newEnt.Url
	}

	if newEnt.Expression != "" {
		e.Expression = newEnt.Expression
	}

	if newEnt.Status != 0 {
		e.Status = newEnt.Status
	}

	if newEnt.Retries > 0 {
		e.Retries = newEnt.Retries
	}

	if newEnt.RequestTimeout > 0 {
		e.RequestTimeout = newEnt.RequestTimeout
	}

	if !newEnt.StopAt.IsEmptyTime() {
		e.StopAt = newEnt.StopAt
	}
}
