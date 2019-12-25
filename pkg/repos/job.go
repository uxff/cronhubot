package repos

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/uxff/cronhubot/pkg/models"
	"github.com/uxff/cronhubot/pkg/utils"
)

type JobRepo interface {
	Create(ent *models.CronJob) (err error)
	FindById(id int) (ent *models.CronJob, err error)
	Update(ent *models.CronJob) (err error)
	Delete(ent *models.CronJob) (err error)
	Search(query *models.Query) (ents []models.CronJob, err error)
}

type CronJob struct {
	db *xorm.Engine
}

func NewCronJob(db *xorm.Engine) *CronJob {
	return &CronJob{db}
}

func (r *CronJob) Create(e *models.CronJob) error {
	e.CreatedAt = utils.NewJsonTimeNow()
	_, err := r.db.Insert(e)
	return err
}

func (r *CronJob) FindById(id int) (e *models.CronJob, err error) {
	e = new(models.CronJob)
	_, err = r.db.Where("id=?", id).Get(e)
	return
}

func (r *CronJob) Update(e *models.CronJob) error {
	e.UpdatedAt = utils.NewJsonTimeNow()
	_, err := r.db.Where("id = ?", e.Id).Update(e)
	return err
}

// 软删除，设置状态为"inactive"即可
func (r *CronJob) Delete(e *models.CronJob) error {
	e.Status = models.StatusInactive
	e.UpdatedAt = utils.NewJsonTimeNow()
	e.StopAt = utils.NewJsonTimeNow()
	_, err := r.db.Where("id = ?", e.Id).Update(e)
	return err
}

func (r *CronJob) Search(q *models.Query) (ents []models.CronJob, err error) {
	if q.IsEmpty() {
		err = r.db.In("stop_at", models.DefaultTimeStr, utils.ZeroTimeStr).Or("stop_at > ?", utils.TimeToString(time.Now())).Where("status = ?", models.StatusActive).Find(&ents)
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

	err = db.Find(ents)

	return
}
