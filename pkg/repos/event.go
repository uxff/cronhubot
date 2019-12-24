package repos

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/uxff/cronhubot/pkg/models"
	"github.com/uxff/cronhubot/pkg/utils"
)

type EventRepo interface {
	Create(event *models.CronJobs) (err error)
	FindById(id int) (event *models.CronJobs, err error)
	Update(event *models.CronJobs) (err error)
	Delete(event *models.CronJobs) (err error)
	Search(query *models.Query) (events []models.CronJobs, err error)
}

type Event struct {
	db *xorm.Engine
}

func NewEvent(db *xorm.Engine) *Event {
	return &Event{db}
}

func (r *Event) Create(e *models.CronJobs) error {
	e.CreatedAt = utils.NewJsonTimeNow()
	_, err := r.db.Insert(e)
	return err
}

func (r *Event) FindById(id int) (e *models.CronJobs, err error) {
	e = new(models.CronJobs)
	_, err = r.db.Where("id=?", id).Get(e)
	return
}

func (r *Event) Update(e *models.CronJobs) error {
	e.UpdatedAt = utils.NewJsonTimeNow()
	_, err := r.db.Where("id = ?", e.Id).Update(e)
	return err
}

// 软删除，设置状态为"inactive"即可
func (r *Event) Delete(e *models.CronJobs) error {
	e.Status = models.Inactive
	e.UpdatedAt = utils.NewJsonTimeNow()
	e.StopAt = utils.NewJsonTimeNow()
	_, err := r.db.Where("id = ?", e.Id).Update(e)
	return err
}

func (r *Event) Search(q *models.Query) (events []models.CronJobs, err error) {
	if q.IsEmpty() {
		err = r.db.In("stop_at", models.DefaultTimeStr, utils.ZeroTimeStr).Or("stop_at > ?", utils.TimeToString(time.Now())).Where("status = ?", models.Active).Find(&events)
		if err != nil {
			return
		}

		return
	}

	var db = r.db
	if q.Status != "" {
		db.Where("status = ?", q.Status)
	}

	if q.Expression != "" {
		db.Where("expression = ?", q.Expression)
	}

	err = db.Find(events)

	return
}
