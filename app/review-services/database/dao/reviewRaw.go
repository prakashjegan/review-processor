package dao

import (
	"fmt"

	mp "github.com/geraldo-labs/merge-struct"
	gdatabase "github.com/prakashjegan/review-processor/app/database"
	"github.com/prakashjegan/review-processor/app/review-services/database/model"
	modeld "github.com/prakashjegan/review-processor/app/review-services/database/model"
	"github.com/prakashjegan/review-processor/app/review-services/utils"
)

type ReviewRawDao interface {
	CreateReviewRaw(reviewRaw *modeld.ReviewRaw) (*modeld.ReviewRaw, error)
	GetReviewRawById(reviewRawById uint64) (reviewRaw modeld.ReviewRaw, err error)
}

func GetReviewRawDao() (reviewRawDao ReviewRawDao) {
	return &reviewRawDaoImpl{}
}

type reviewRawDaoImpl struct {
}

func (r *reviewRawDaoImpl) CreateReviewRaw(reviewRaw *modeld.ReviewRaw) (*modeld.ReviewRaw, error) {
	if reviewRaw.RawIdentifier == "" {
		reviewRaw.RawIdentifier = fmt.Sprintf("%d", reviewRaw.ID)
	}
	db := gdatabase.GetDB()
	oldReviewRaw := &model.ReviewRaw{}
	er := db.Where("raw_identifier = ? ", reviewRaw.RawIdentifier).First(oldReviewRaw).Error
	if er != nil {
		oldReviewRaw.ID = utils.GetUID()
	}
	mp.Struct(oldReviewRaw, reviewRaw)
	er = db.Save(oldReviewRaw).Error
	return oldReviewRaw, er
}

func (r *reviewRawDaoImpl) GetReviewRawById(reviewRawId uint64) (reviewRaw modeld.ReviewRaw, err error) {
	db := gdatabase.GetDB()
	err = db.Where("id = ? ", reviewRawId).First(&reviewRaw, reviewRawId).Error
	return
}
