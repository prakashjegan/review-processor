package dao

// import (
// 	model "github.com/review-processor/app/review-services/database/model"
// )

// func GetJobEventsDao() (jobEventDao JobEventDao)

// type JobEventDao interface {
// 	GetJobEvents(string) (events []model.JobEvent, err error)
// }

// type jobEventDaoImpl struct {
// }

// func (jed *jobEventDaoImpl) GetJobEvents(jobType string) (events []model.JobEvent, err error) {
// 	db := database.GetDB()
// 	err = db.Where("job_type=?", jobType).Order("id desc").Limit(50).Find(&events).Error
// 	return events, err
// }
