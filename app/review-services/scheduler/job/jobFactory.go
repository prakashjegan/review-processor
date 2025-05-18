package scheduler

import (
	"errors"

	"github.com/prakashjegan/review-processor/app/config"
	"github.com/reugn/go-quartz/quartz"
)

func NewJobFactory() (jobFactory JobFactory) {
	return &JobFactoryImpl{}
}

type JobFactory interface {
	CreateJob(jobType string) (job quartz.Job, err error)
}

type JobFactoryImpl struct {
	Jobs map[string]quartz.Job
}

func (j *JobFactoryImpl) CreateJob(jobType string) (job quartz.Job, err error) {
	if j.Jobs == nil {
		j.Jobs = make(map[string]quartz.Job)
	}
	if job, ok := j.Jobs[jobType]; ok {
		return job, nil
	}

	switch jobType {
	case "REVIEW_SCHEDULER":
		// Get configuration from the config package
		config := config.GetConfig()
		job = NewReviewJob(config)
	default:
		err = errors.New("unknown job type")
	}
	j.Jobs["REVIEW_SCHEDULER"] = job
	return job, err
}
