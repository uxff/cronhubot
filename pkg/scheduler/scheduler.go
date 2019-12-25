package scheduler

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/robfig/cron"
	"github.com/uxff/cronhubot/pkg/models"
	"github.com/uxff/cronhubot/pkg/repos"
	"github.com/uxff/cronhubot/pkg/runner"
	"github.com/uxff/cronhubot/pkg/log"
)

type Scheduler interface {
	Create(event *models.CronJob) error
	Update(event *models.CronJob) error
	Delete(id uint) error
	Find(id uint) (*cron.Cron, error)
	ScheduleAll()
}

type scheduler struct {
	sync.RWMutex
	Kv map[uint]*cron.Cron
	// Cron *cron.Cron
	r repos.JobRepo
}

func New(r repos.JobRepo) Scheduler {
	s := &scheduler{
		Kv: make(map[uint]*cron.Cron),
		// Cron: cron.New(),
		r: r,
	}

	// s.Cron.Start()

	return s
}

func (s *scheduler) Create(e *models.CronJob) (err error) {
	newCron := cron.New()

	// 给回调地址上添加定时任务id
	e.Url = UrlAddParam(e.Url, map[string]string{"__cronId": strconv.Itoa(int(e.Id))})

	err = newCron.AddFunc(e.Expression, func() {
		defer func() {
			log.Debugf("结束定时函数, id:%d", e.Id)
		}()
		traceId := fmt.Sprintf("%d", e.Id)
		// 如果定时任务已过期，不执行
		// 1、删除内存里的定时
		// 2、更改数据库里的状态
		if !e.StopAt.IsEmptyTime() && e.StopAt.ToStdTime().Before(time.Now()) {
			go func() { // 删除需要放在协程里，否则可能会死锁
				log.Trace(traceId).Infof("当前时间已经大于定时任务的结束时间，准备删除数据, 任务配置%+v", e)
				if err := s.Delete(e.Id); err != nil {
					log.Trace(traceId).Errorf("从内存删除定时任务失败:%v", err)
				}
				if err := s.r.Delete(e); err != nil {
					log.Trace(traceId).Errorf("从DB删除定时任务失败:%v", err)
				}
			}()
			return
		}

		c := &runner.Config{
			Url:     e.Url,
			Retries: e.Retries,
			Timeout: e.RequestTimeout,
		}

		resp, err := runner.NoticeRemote(traceId, c)
		if err != nil {
			log.Trace(traceId).Errorf("通知业务方失败：%v 配置:%+v", err, c)
			return
		}

		// 如果远端回复定时任务不需要继续执行，则删除定时任务
		if resp.UselessJob() {
			go func() { // 删除需要放在协程里，否则可能会死锁
				if err := s.Delete(e.Id); err != nil {
					log.Trace(traceId).Errorf("从内存删除定时任务失败:%v", err)
				}
				if err := s.r.Delete(e); err != nil {
					log.Trace(traceId).Errorf("从DB删除定时任务失败:%v", err)
				}
			}()
		}
	})

	if err != nil {
		log.Errorf("定时任务执行失败：%v id:%d", err, e.Id)
		return err
	}

	s.Lock()
	defer s.Unlock()

	newCron.Start()
	s.Kv[e.Id] = newCron

	return
}

func (s *scheduler) Find(id uint) (cronJob *cron.Cron, err error) {
	s.Lock()
	defer s.Unlock()

	cronJob, found := s.Kv[id]
	if !found {
		err = ErrEventNotExist
		return
	}

	return
}

func (s *scheduler) Update(e *models.CronJob) (err error) {
	if err = s.Delete(e.Id); err != nil {
		return
	}

	return s.Create(e)
}

func (s *scheduler) Delete(id uint) (err error) {
	s.Lock()
	defer s.Unlock()

	_, found := s.Kv[id]
	if !found {
		log.Errorf("定时任务未找到,id = %d", id)
		err = ErrEventNotExist
		return
	}

	s.Kv[id].Stop()
	s.Kv[id] = nil
	log.Infof("定时任务删除成功,id = %d", id)

	// for k := range s.Kv {
	// 	log.Infof("最新的定时任务：%V", k)
	// }
	return
}

func (s *scheduler) ScheduleAll() {
	events, err := s.r.Search(&models.Query{})
	if err != nil {
		log.Errorf("Failed to find events!")
		return
	}

	for index, e := range events {
		log.Infof("准备启动定时任务(%d)：%+v", index, e)
		if err = s.Create(&events[index]); err != nil {
			log.Errorf("Failed to create event! event:%+v", e)
		}
	}
}

func UrlAddParam(targetUrl string, params map[string]string) string {
	urlParsed, err := url.Parse(targetUrl)
	if err != nil {
		return targetUrl
	}

	u := urlParsed.Query()

	for paramKey, paramVal := range params {
		u.Set(paramKey, paramVal)
	}
	urlParsed.RawQuery = u.Encode()

	return urlParsed.String()
}
