package processor

import (
	"github.com/prakashjegan/review-processor/app/review-services/database/model"
	models "github.com/prakashjegan/review-processor/app/review-services/models"
	pro "github.com/prakashjegan/review-processor/app/review-services/redisstream/processor/transformer"
)

type Processor interface {
	Transform(model.ReviewRaw) (*models.HotelReviewData, error)
}

func NewProcessor(providerType string) (Processor, error) {
	switch providerType {
	case "hotel":
		return pro.GetDefaultProcessor(), nil
	default:
		return pro.GetDefaultProcessor(), nil
	}
}
