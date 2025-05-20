package processor

import (
	"encoding/json"
	"fmt"

	gdatabase "github.com/prakashjegan/review-processor/app/database"
	"github.com/prakashjegan/review-processor/app/review-services/database/dao"
	"github.com/prakashjegan/review-processor/app/review-services/database/model"
	models "github.com/prakashjegan/review-processor/app/review-services/models"
	"github.com/prakashjegan/review-processor/app/review-services/redisstream"
	log "github.com/sirupsen/logrus"
)

func ProcessReview(message redisstream.Message) {

	data := message.Data
	rawData := data["raw_data"]
	rawDataS := fmt.Sprintf("%v", rawData)
	reviewRaw := &model.ReviewRaw{}
	err := json.Unmarshal([]byte(rawDataS), reviewRaw)
	if err != nil {
		log.Info("Error in unmarshalling raw data")
		return
	}
	fmt.Printf("Received Review: %s\n", *reviewRaw)

	processor, err := NewProcessor(reviewRaw.ThirdPartyProviderName)
	rawDataDao := dao.GetReviewRawDao()
	if err != nil {
		updateFailure(reviewRaw, rawDataDao, err.Error())

		log.Info(err.Error())
		return
	}
	hotelReview, err := processor.Transform(*reviewRaw)
	if err != nil {
		updateFailure(reviewRaw, rawDataDao, err.Error())
		log.Info(err.Error())
		return
	}

	if hotelReview == nil {
		updateFailure(reviewRaw, rawDataDao, "Parsing Failure")
		return
	}

	db := gdatabase.GetDB()
	tx := db.Begin()
	// Populate ProductReview
	productReviewDAO := dao.NewProductReviewDAO(tx)
	productReview := &model.ProductReview{
		HotelID:   hotelReview.HotelID,
		HotelName: hotelReview.HotelName,
	}
	var err1 error
	productReview, err1 = productReviewDAO.CreateOrUpdate(productReview)
	if err1 != nil {
		tx.Rollback()
		updateFailure(reviewRaw, rawDataDao, fmt.Sprintf("Error while creating or updating ProductReview : %s", err1.Error()))
		log.Info(err1.Error())
		return
	}

	//Populate Comment
	comment := hotelReview.Comment
	reviewCommentDAO := dao.NewReviewCommentDAO(tx)
	reviewComment := &model.ReviewComment{
		ProductReviewId:       productReview.ID,
		HotelID:               hotelReview.HotelID,
		HotelName:             hotelReview.HotelName,
		IsShowReviewResponse:  comment.IsShowReviewResponse,
		Rating:                comment.Rating,
		HotelReviewID:         comment.HotelReviewID,
		ProviderID:            comment.ProviderID,
		CheckInDateMonthYear:  comment.CheckInDateMonthYear,
		EncryptedReviewData:   comment.EncryptedReviewData,
		FormattedRating:       comment.FormattedRating,
		FormattedReviewDate:   comment.FormattedReviewDate,
		FormattedResponseDate: comment.FormattedResponseDate,
		RatingText:            comment.RatingText,
		ResponderName:         comment.ResponderName,
		ResponseDateText:      comment.ResponseDateText,
		ResponseTranslateSrc:  comment.ResponseTranslateSrc,

		ReviewComments:     comment.ReviewComments,
		ReviewNegatives:    comment.ReviewNegatives,
		ReviewPositives:    comment.ReviewPositives,
		ReviewProviderLogo: comment.ReviewProviderLogo,
		ReviewProviderText: comment.ReviewProviderText,
		ReviewTitle:        comment.ReviewTitle,
		TranslateSource:    comment.TranslateSource,
		TranslateTarget:    comment.TranslateTarget,
		ReviewDate:         comment.ReviewDate,
		OriginalTitle:      comment.OriginalTitle,
		OriginalComment:    comment.OriginalComment,
	}
	reviewComment, err1 = reviewCommentDAO.CreateOrUpdate(reviewComment)
	if err1 != nil {
		tx.Rollback()
		updateFailure(reviewRaw, rawDataDao, fmt.Sprintf("Error while creating or updating Review comment : %s", err1.Error()))
		log.Info(err1.Error())
		return
	}

	reviewerInfoDAO := dao.NewReviewerInfoDAO(tx)
	reviewer := comment.ReviewerInfo
	reviewerInfo := &model.ReviewerInfo{

		ReviewCommentId: reviewComment.ID,
		ProductReviewId: reviewComment.ProductReviewId,

		CountryName:           reviewer.CountryName,
		DisplayMemberName:     reviewer.DisplayMemberName,
		FlagName:              reviewer.FlagName,
		ReviewGroupName:       reviewer.ReviewGroupName,
		RoomTypeName:          reviewer.RoomTypeName,
		CountryID:             reviewer.CountryID,
		LengthOfStay:          reviewer.LengthOfStay,
		ReviewGroupID:         reviewer.ReviewGroupID,
		RoomTypeID:            reviewer.RoomTypeID,
		ReviewerReviewedCount: reviewer.ReviewerReviewedCount,
		IsExpertReviewer:      reviewer.IsExpertReviewer,
		IsShowGlobalIcon:      reviewer.IsShowGlobalIcon,
		IsShowReviewedCount:   reviewer.IsShowReviewedCount,
	}

	reviewerInfo, err1 = reviewerInfoDAO.CreateOrUpdate(reviewerInfo)
	if err1 != nil {
		tx.Rollback()
		updateFailure(reviewRaw, rawDataDao, fmt.Sprintf("Error while creating or updating Reviewer Info : %s", err1.Error()))
		log.Info(err1.Error())
		return
	}

	providers := hotelReview.OverallByProviders
	if len(providers) > 0 {
		productReviewByProviderDAO := dao.NewProductReviewByProviderDAO(tx)
		productReviewGradeDAO := dao.NewProductReviewGradeDAO(tx)
		for _, provider := range providers {

			productReviewByProvider := &model.ProductReviewByProvider{
				ProductReviewId: productReview.ID,
				HotelID:         productReview.HotelID,
				HotelName:       productReview.HotelName,
				Platform:        provider.Provider,
				ProviderId:      provider.ProviderID,
				Provider:        provider.Provider,
				OverallScore:    int64(provider.OverallScore * 100),
				ReviewCount:     int64(provider.ReviewCount),
			}
			productReviewByProvider, err1 = productReviewByProviderDAO.CreateOrUpdate(productReviewByProvider)
			if err1 != nil {
				tx.Rollback()
				updateFailure(reviewRaw, rawDataDao, fmt.Sprintf("Error while creating or updating ProductReview By Provider : %s", err1.Error()))
				log.Info(err1.Error())
				return
			}
			grades := provider.Grades
			if len(grades) <= 0 {
				continue
			}
			for attribute, grade := range grades {
				productReviewGrade := &model.ProductReviewGrade{
					ProductReviewByProviderId: productReviewByProvider.ID,
					ProductReviewId:           productReview.ID,
					HotelID:                   productReview.HotelID,
					HotelName:                 productReview.HotelName,
					Platform:                  productReviewByProvider.Platform,
					ProviderId:                productReviewByProvider.ProviderId,
					Provider:                  productReviewByProvider.Provider,
					ThirdPartyProviderId:      productReviewByProvider.ThirdPartyProviderId,
					GradeAttribute:            attribute,
					Grade:                     int64(grade * 100),
				}
				productReviewGrade, err1 = productReviewGradeDAO.CreateOrUpdate(productReviewGrade)
				if err1 != nil {
					tx.Rollback()
					updateFailure(reviewRaw, rawDataDao, fmt.Sprintf("Error while creating or updating ProductReview Grade : %s", err1.Error()))
					log.Info(err1.Error())
					return
				}
			}
		}
	}

	tx.Commit()

	CalculateMetricsForProductReview(productReview)
	//TODO : Calculate the over All performance of the hotel and update it to database.

	updateSuccess(reviewRaw, rawDataDao)
	log.Infof("\nSuccessfully processed review for Hotel ID : %d", hotelReview.HotelID)

}

func updateFailure(reviewRaw *model.ReviewRaw, rawDataDao dao.ReviewRawDao, message string) {
	reviewRaw.Status = string(models.FileStatusFailed)
	reviewRaw.Message = message
	_, e := rawDataDao.CreateReviewRaw(reviewRaw)
	if e != nil {
		log.Info(e.Error())
		return
	}
	return
}

func updateSuccess(reviewRaw *model.ReviewRaw, rawDataDao dao.ReviewRawDao) {
	reviewRaw.Status = string(models.JobStatusCompleted)
	reviewRaw.Message = "success"
	_, e := rawDataDao.CreateReviewRaw(reviewRaw)
	if e != nil {
		log.Info(e.Error())
		return
	}
	return
}

func CalculateMetricsForProductReview(productReview *model.ProductReview) (*model.ProductReview, error) {
	db := gdatabase.GetDB()
	providerDao := dao.NewProductReviewByProviderDAO(db)
	metrics, err := providerDao.GetMetricsByProductId(productReview.ID)
	if err != nil {
		return nil, err
	}
	if len(metrics) <= 0 { //No metrics available.
		return productReview, nil
	}
	metric := metrics[0]
	productReviewDao := dao.NewProductReviewDAO(db)
	productReview.OverallScore = int64(metric.OverAllPerformance)
	productReview.ReviewCount = metric.TotalReviewCount
	productReview, err = productReviewDao.CreateOrUpdate(productReview)
	return productReview, nil
}
