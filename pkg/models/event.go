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
	Active         = 1
	Inactive       = 9
	DefaultTimeStr = "1999-01-01 00:00:00"
)

type CronJobs struct {
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

func NewEvent() *CronJobs {
	return &CronJobs{
		Status:         Active, // 默认是有效的定时任务
		Retries:        0,      // 默认请求业务方时不重试
		RequestTimeout: 3,      // 请求业务方时默认超时时间是3s
	}
}

func (e *CronJobs) TableName() string {
	return "athena_cronjobs"
}

func (e *CronJobs) CheckExpression(expression string) error {
	if len(strings.Split(expression, " ")) != 6 {
		return errors.New("定时配置格式不正确")
	}

	if _, err := cron.Parse(expression); err != nil {
		return err
	}
	return nil
}

func (e *CronJobs) Validate() (errors map[string]string, ok bool) {
	errors = make(map[string]string)
	if e.Url == "" {
		errors["url"] = "url is empty"
	}

	if err := e.CheckExpression(e.Expression); err != nil {
		errors["expression"] = err.Error()
	}

	if e.Status != Active && e.Status != Inactive {
		errors["status"] = fmt.Sprintf("status must be %v or %v", Active, Inactive)
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

func (e *CronJobs) SetAttributes(newEvent *CronJobs) {
	if newEvent.Name != "" {
		e.Name = newEvent.Name
	}

	if newEvent.Url != "" {
		e.Url = newEvent.Url
	}

	if newEvent.Expression != "" {
		e.Expression = newEvent.Expression
	}

	if newEvent.Status != 0 {
		e.Status = newEvent.Status
	}

	if newEvent.Retries > 0 {
		e.Retries = newEvent.Retries
	}

	if newEvent.RequestTimeout > 0 {
		e.RequestTimeout = newEvent.RequestTimeout
	}

	if !newEvent.StopAt.IsEmptyTime() {
		e.StopAt = newEvent.StopAt
	}
}
