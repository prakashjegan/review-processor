package dao

import (
	mp "github.com/geraldo-labs/merge-struct"
	gdatabase "github.com/prakashjegan/review-processor/app/database"
	modeld "github.com/prakashjegan/review-processor/app/review-services/database/model"
	rutils "github.com/prakashjegan/review-processor/app/review-services/utils"
)

type ReviewFileStatesDao interface {
	CreateReviewFileStates(reviewFileStates *modeld.ReviewFileStates) (*modeld.ReviewFileStates, error)
	GetReviewFileStates(fileName string) (reviewFileStates modeld.ReviewFileStates, err error)
	GetLastProcessedFile(thirdParty string) (reviewFileStates modeld.ReviewFileStates, err error)
}

func GetReviewFileStatesDao() (reviewFileStatesDao ReviewFileStatesDao) {
	return &reviewFileStatesDaoImpl{}
}

type reviewFileStatesDaoImpl struct {
}

func (rfs *reviewFileStatesDaoImpl) CreateReviewFileStates(reviewFileStates *modeld.ReviewFileStates) (*modeld.ReviewFileStates, error) {

	db := gdatabase.GetDB()
	oldReviewFileStates := modeld.ReviewFileStates{}
	if db.Where("check_sum=?", reviewFileStates.CheckSum).First(&oldReviewFileStates).Error != nil {
		oldReviewFileStates.Id = rutils.GetUID()
	}
	mp.Struct(&oldReviewFileStates, reviewFileStates)

	err := db.Save(&oldReviewFileStates).Error
	if err != nil {
		return nil, err
	}
	return &oldReviewFileStates, nil

}

func (rfs *reviewFileStatesDaoImpl) GetLastProcessedFile(thirdParty string) (reviewFileStates modeld.ReviewFileStates, err error) {
	db := gdatabase.GetDB()
	if err = db.Where("third_party_name=? and state='SUCCESS'", thirdParty).Order(" file_id desc").Limit(1).Find(&reviewFileStates).Error; err != nil {
		return
	}
	return
}

func (rfs *reviewFileStatesDaoImpl) GetReviewFileStates(checksum string) (reviewFileStates modeld.ReviewFileStates, err error) {
	db := gdatabase.GetDB()
	if err = db.Where("check_sum=?", checksum).Find(&reviewFileStates).Error; err != nil {
		return
	}
	return
}
