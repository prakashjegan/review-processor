package dao

import (
	mp "github.com/geraldo-labs/merge-struct"
	gdatabase "github.com/prakashjegan/review-processor/app/database"
	modeld "github.com/prakashjegan/review-processor/app/review-services/database/model"
	"github.com/prakashjegan/review-processor/app/review-services/utils"
)

func GetJobEventsDao() (jobEventDao JobEventDao)

type JobEventDao interface {
	GetJobEvents(string) (events []modeld.JobEvent, err error)
	CreateOrUpdateJobEvent(jobEvent *modeld.JobEvent) (newJobEvent *modeld.JobEvent, err error)
}

func GetJobEventDao() (jobEventDao JobEventDao) {
	return &jobEventDaoImpl{}
}

type jobEventDaoImpl struct {
}

func (jed *jobEventDaoImpl) GetJobEvents(jobType string) (events []modeld.JobEvent, err error) {
	db := gdatabase.GetDB()
	err = db.Where("job_type=?", jobType).Order("id desc").Limit(50).Find(&events).Error
	return events, err
}

func (jed *jobEventDaoImpl) CreateOrUpdateJobEvent(jobEvent *modeld.JobEvent) (newJobEvent *modeld.JobEvent, err error) {
	db := gdatabase.GetDB()
	oldjobEvent := &modeld.JobEvent{}
	if jobEvent.Id != 0 {
		db.Where("job_id=? and event_name=?", jobEvent.JobId, jobEvent.EventName).First(oldjobEvent)
	}
	if oldjobEvent.Id == 0 {
		oldjobEvent.Id = utils.GetUID()
	}
	mp.Struct(oldjobEvent, jobEvent)

	err = db.Save(oldjobEvent).Error
	return oldjobEvent, err
}
