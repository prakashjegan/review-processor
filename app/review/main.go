package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	gconfig "github.com/prakashjegan/review-processor/app/config"
	gdatabase "github.com/prakashjegan/review-processor/app/database"
	"github.com/prakashjegan/review-processor/app/database/migrate"
	sch "github.com/prakashjegan/review-processor/app/review-services/scheduler"
	rjob "github.com/prakashjegan/review-processor/app/review-services/scheduler/job"
)

func main() {
	// Load configuration
	err := gconfig.Config()
	if err != nil {
		log.Info(err)
		return
	}
	log.Info("Starting up...")
	// read configs
	configure := gconfig.GetConfig()
	log.Info("Configurations :::: %+v :::: %s ", configure, configure.Database.RDBMS.Activate)
	if configure.Database.RDBMS.Activate == gconfig.Activated {
		// Initialize RDBMS client
		log.Info("Configurations :::: %+v :::: %s ", configure, configure.Database.RDBMS.Activate)

		if err := gdatabase.InitDB().Error; err != nil {
			log.Info(err)
			return
		}

		// TODO ::::: Use Only this code.
		if err := migrate.DropAllTablesWithActual(); err != nil {
			log.Info(err)
			return
		}

		// Start DB migration
		if err := migrate.StartMigrationWithActualData(*configure); err != nil {
			log.Info(err)
			return
		}
	}

	if configure.Database.REDIS.Activate == gconfig.Activated {
		// Initialize REDIS client
		if _, err := gdatabase.InitRedis(); err != nil {
			log.Info(err)
			return
		}
	}

	scheduler := sch.NewReviewScheduler(configure)
	scheduler.Start()

	job := rjob.NewReviewJob(configure)
	job.Execute(context.Background())

	select {}
	log.Println("Shutting down...")
}
