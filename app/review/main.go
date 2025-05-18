package main

import (
	"fmt"
	"log"

	gconfig "github.com/prakashjegan/review-processor/app/config"
	gdatabase "github.com/prakashjegan/review-processor/app/database"
	"github.com/prakashjegan/review-processor/app/database/migrate"
	sch "github.com/prakashjegan/review-processor/app/review-services/scheduler"
)

func main() {
	// Load configuration
	err := gconfig.Config()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Completed Config")
	// read configs
	configure := gconfig.GetConfig()

	if configure.Database.RDBMS.Activate == gconfig.Activated {
		// Initialize RDBMS client
		if err := gdatabase.InitDB().Error; err != nil {
			fmt.Println(err)
			return
		}

		// TODO ::::: Use Only this code.
		if err := migrate.DropAllTablesWithActual(); err != nil {
			fmt.Println(err)
			return
		}

		// Start DB migration
		if err := migrate.StartMigrationWithActualData(*configure); err != nil {
			fmt.Println(err)
			return
		}
	}

	if configure.Database.REDIS.Activate == gconfig.Activated {
		// Initialize REDIS client
		if _, err := gdatabase.InitRedis(); err != nil {
			fmt.Println(err)
			return
		}
	}

	scheduler := sch.NewReviewScheduler(configure)
	scheduler.Start()

	log.Println("Shutting down...")
}
