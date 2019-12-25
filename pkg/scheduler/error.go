package scheduler

import (
	"errors"
)

var (
	ErrCronjobNotExist = errors.New("finding a scheduled cronjob requires a existent cron id")
)
