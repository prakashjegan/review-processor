package processor

import (
	"encoding/json"

	"github.com/prakashjegan/review-processor/app/review-services/database/model"
	"github.com/prakashjegan/review-processor/app/review-services/models"
	log "github.com/sirupsen/logrus"
)

type DefaultProcessor struct {
}

func GetDefaultProcessor() *DefaultProcessor {
	return &DefaultProcessor{}
}

func (d *DefaultProcessor) Transform(data model.ReviewRaw) (*models.HotelReviewData, error) {
	log.Infof("\nTransforming the raw data :: %d", data.ID)
	dat := data.RawData
	reviewData := &models.HotelReviewData{}
	er := json.Unmarshal([]byte(dat), reviewData)
	if er != nil {
		log.Errorf("Error while transforming the data :: %s", er.Error())
		return nil, er
	}
	log.Infof("\nSuccessfully transformed the data :: %+v", reviewData)
	return reviewData, nil
}
