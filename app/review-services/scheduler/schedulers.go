package scheduler

import (
	"context"
	"fmt"

	"github.com/prakashjegan/review-processor/app/config"
	job "github.com/prakashjegan/review-processor/app/review-services/scheduler/job"
	"github.com/reugn/go-quartz/quartz"
)

type ReviewScheduler struct {
	schedulers map[string]quartz.Scheduler
	Config     *config.Configuration
}

func NewReviewScheduler(config *config.Configuration) *ReviewScheduler {
	return &ReviewScheduler{
		Config:     config,
		schedulers: make(map[string]quartz.Scheduler),
	}
}

func (s *ReviewScheduler) Start() error {

	// Create a cron trigger
	config := s.Config.Scheduler
	for name, triggerV := range config.SchedulerCrons {
		if _, ok := s.schedulers[name]; !ok {
			sq := quartz.NewStdScheduler()

			s.schedulers[name] = sq
		}
		scheduler := s.schedulers[name]
		trigger, err := quartz.NewCronTrigger(triggerV) // Run every 5 minutes
		if err != nil {
			return fmt.Errorf("failed to create trigger: %v", err)
		}
		// Schedule the job
		job, err := job.NewJobFactory().CreateJob(name)
		if err != nil {
			return fmt.Errorf("failed to create job: %v", err)
		}
		ctx := context.Background()
		err = scheduler.ScheduleJob(ctx, job, trigger)
		if err != nil {
			return fmt.Errorf("failed to schedule job: %v", err)
		}
	}

	return nil
}

func (s *ReviewScheduler) Stop() error {
	config := s.Config.Scheduler
	for cron, _ := range config.SchedulerCrons {
		if _, ok := s.schedulers[cron]; !ok {
			sq := quartz.NewStdScheduler()

			s.schedulers[cron] = sq
		}
		scheduler := s.schedulers[cron]
		// Schedule the job
		scheduler.Stop()

	}
	return nil
}
