package datastore

import (
	"net/url"
	"strings"

	"github.com/go-xorm/xorm"
)

const (
	MySQL = "mysql"
)

func New(dsn string) (*xorm.Engine, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case MySQL:
		c := MysqlConfig{
			Url:         strings.Replace(dsn, "mysql://", "", 1),
			MaxIdleConn: 10,
			MaxOpenConn: 100,
			LogMode:     true,
		}
		return NewMysql(c)
	default:
		return nil, ErrUnknownDatabaseProvider
	}
}
