// Package migrate to migrate the schema
package migrate

import (
	"fmt"

	gconfig "github.com/prakashjegan/review-processor/app/config"
	gdatabase "github.com/prakashjegan/review-processor/app/database"
	modeld "github.com/prakashjegan/review-processor/app/review-services/database/model"
)

type jobEvent modeld.JobEvent
type productReview modeld.ProductReview
type reviewComment modeld.ReviewComment
type reviewerInfo modeld.ReviewerInfo
type productReviewByProvider modeld.ProductReviewByProvider
type productReviewGrade modeld.ProductReviewGrade
type reviewFileStates modeld.ReviewFileStates
type reviewRaw modeld.ReviewRaw
type thirdPartyConfig modeld.ThirdPartyConfig
type transformerMapping modeld.TransformerMapping
type thirdPartyEvent modeld.ThirdPartyEvent

type productReviewEvent modeld.ProductReviewEvent
type reviewCommentEvent modeld.ReviewCommentEvent
type reviewerInfoEvent modeld.ReviewerInfoEvent
type productReviewByProviderEvent modeld.ProductReviewByProviderEvent
type productReviewGradeEvent modeld.ProductReviewGradeEvent

func ShouldMigrateLimited() bool {

	return true
	//return true
}

// DropAllTables - careful! It will drop all the tables!
func DropAllTablesWithActual() error {
	//doMigration := false
	doMigration := ShouldMigrateLimited()
	if !doMigration {
		return nil
	}
	db := gdatabase.GetDB()

	if err := db.Migrator().DropTable(
		&jobEvent{},
		&productReview{},
		&reviewComment{},
		&reviewerInfo{},
		&productReviewByProvider{},
		&productReviewGrade{},
		&reviewFileStates{},
		&reviewRaw{},
		&thirdPartyConfig{},
		&transformerMapping{},
		&thirdPartyEvent{},
		&productReviewEvent{},
		&reviewCommentEvent{},
		&reviewerInfoEvent{},
		&productReviewByProviderEvent{},
		&productReviewGradeEvent{},
	); err != nil {
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

	if err := db.AutoMigrate(
		&jobEvent{},
		&productReview{},
		&reviewComment{},
		&reviewerInfo{},
		&productReviewByProvider{},
		&productReviewGrade{},
		&reviewFileStates{},
		&reviewRaw{},
		&thirdPartyConfig{},
		&transformerMapping{},
		&thirdPartyEvent{},
		&productReviewEvent{},
		&reviewCommentEvent{},
		&reviewerInfoEvent{},
		&productReviewByProviderEvent{},
		&productReviewGradeEvent{},
	); err != nil {

		return err
	}

	fmt.Println("new tables are  migrated successfully!")
	return nil
}
