package checker

import (
	"strings"

	"github.com/go-xorm/xorm"
)

type mysql struct {
	url string
}

func NewMysql(url string) *mysql {
	return &mysql{url: strings.Replace(url, "mysql://", "", 1)}
}

func (p *mysql) IsAlive() bool {
	conn, err := xorm.NewEngine("mysql", p.url)
	if err != nil {
		return false
	}
	defer conn.Close()

	if err = conn.DB().Ping(); err != nil {
		return false
	}

	return true
}
