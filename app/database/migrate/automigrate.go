// Package migrate to migrate the schema
package migrate

import (
	"fmt"

	gconfig "github.com/prakashjegan/review-processor/app/config"
	gdatabase "github.com/prakashjegan/review-processor/app/database"
)

// DropAllTables - careful! It will drop all the tables!

func ShouldMigrate() bool {
	return false
	//return true
}

func ShouldMigrateLimited() bool {

	return false
	//return true
}
func DropAllTables() error {
	db := gdatabase.GetDB()

	if err := db.Migrator().DropTable(); err != nil {
		return err
	}

	fmt.Println("old tables are deleted!")
	return nil
}

// DropAllTables - careful! It will drop all the tables!
func DropAllTablesLimited() error {
	//doMigration := false
	doMigration := ShouldMigrate()
	if !doMigration {
		return nil
	}
	db := gdatabase.GetDB()

	if err := db.Migrator().DropTable(); err != nil {
		return err
	}

	fmt.Println("old tables are deleted!")
	return nil
}

// DropAllTables - careful! It will drop all the tables!
func DropAllTablesWithActual() error {
	//doMigration := false
	doMigration := ShouldMigrateLimited()
	if !doMigration {
		return nil
	}
	db := gdatabase.GetDB()

	if err := db.Migrator().DropTable(); err != nil {
		return err
	}

	fmt.Println("old tables are deleted!")
	return nil
}

func StartMigrationWithActualData(configure gconfig.Configuration) error {
	//doMigration := false
	doMigration := ShouldMigrateLimited()
	if !doMigration {
		return nil
	}
	db := gdatabase.GetDB()
	configureDB := configure.Database.RDBMS
	driver := configureDB.Env.Driver

	if driver == "mysql" {
		// db.Set() --> add table suffix during auto migration
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(); err != nil {
			return err
		}

		fmt.Println("new tables are  migrated successfully!")
		return nil
	}

	if err := db.AutoMigrate(); err != nil {

		return err
	}

	fmt.Println("new tables are  migrated successfully!")
	return nil
}

// StartMigration - automatically migrate all the tables
// - Only create tables with missing columns and missing indexes
// - Will not change/delete any existing columns and their types
func StartMigrationLimited(configure gconfig.Configuration) error {
	//doMigration := false
	doMigration := ShouldMigrate()
	if !doMigration {
		return nil
	}
	db := gdatabase.GetDB()

	if err := db.AutoMigrate(); err != nil {
		return err
	}

	fmt.Println("new tables are  migrated successfully!")
	return nil
}

// StartMigration - automatically migrate all the tables
// - Only create tables with missing columns and missing indexes
// - Will not change/delete any existing columns and their types
func StartMigration(configure gconfig.Configuration) error {
	db := gdatabase.GetDB()

	if err := db.AutoMigrate(); err != nil {
		return err
	}

	fmt.Println("new tables are  migrated successfully!")
	return nil
}
