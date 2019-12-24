package main

import (
	"fmt"
	"time"

	"github.com/uxff/cronhubot/pkg/utils"
)

type AAA struct {
	Created utils.JsonTime
}

type CronJobs struct {
	Id             uint           `json:"id"               gorm:"AUTO_INCREMENT;primary_key;"`
	Name           string         `json:"name"             gorm:"type:varchar(50)"`
	Url            string         `json:"url"              gorm:"type:varchar(256)"`
	Expression     string         `json:"expression"       gorm:"type:varchar(20)"`
	Status         string         `json:"status"           gorm:"type:varchar(10);index:idx_status_expire"`
	Retries        uint           `json:"retries"          gorm:"default:0"`
	RequestTimeout uint           `json:"request_timeout"  gorm:"default:3"`
	ExpireTime     utils.JsonTime `json:"expire_time"      gorm:"type:timestamp;index:idx_status_expire;default:'1900-01-01 00:00:00'"`
	CreatedAt      utils.JsonTime `json:"created_at"       gorm:"type:timestamp;default:'1900-01-01 00:00:00'"`
	UpdatedAt      utils.JsonTime `json:"updated_at"       gorm:"type:timestamp;default:'1900-01-01 00:00:00'"`
}

func Show(o *CronJobs){
	time.Sleep(time.Second*3)
	fmt.Println(o.Id)
}

func main() {
	list := []CronJobs{
		{
			Id:             11,
		},
		{
			Id:             22,
		},
		{
			Id:             33,
		},
	}

	for index := range list{
		go Show(&list[index])
	}


	time.Sleep(time.Hour)
}
