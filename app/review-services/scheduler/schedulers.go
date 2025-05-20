package scheduler

import (
	"context"
	"fmt"

	"github.com/prakashjegan/review-processor/app/config"
	job "github.com/prakashjegan/review-processor/app/review-services/scheduler/job"
	"github.com/reugn/go-quartz/quartz"
	log "github.com/sirupsen/logrus"
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
		log.Info("Starting review processor scheduler : ", name)
		if _, ok := s.schedulers[name]; !ok {
			sq := quartz.NewStdScheduler()

			s.schedulers[name] = sq
			sq.Start(context.Background())
		}
		scheduler := s.schedulers[name]
		trigger, err := quartz.NewCronTrigger(triggerV) // Run every 5 minutes
		log.Info("review processor scheduler Trigger created : ", name)

		if err != nil {
			log.Info("Failed review processor scheduler :", err)
			return fmt.Errorf("failed to create trigger: %v", err)
		}
		// Schedule the job
		job, err := job.NewJobFactory().CreateJob(name)
		log.Info("review processor scheduler Job created : ", name)
		if err != nil {
			log.Info("Failed to create job  :", err)
			return fmt.Errorf("failed to create job: %v", err)
		}
		ctx := context.Background()
		log.Info("starting review processor schedule Job   : ", name)
		err = scheduler.ScheduleJob(ctx, job, trigger)
		log.Info("scheduled review processor scheduler Job   : ", name)
		if err != nil {
			log.Info("Failed to schedule job  :", err)
			return fmt.Errorf("failed to schedule job: %v", err)
		}
		log.Info("Started review processor scheduler :", name)

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
