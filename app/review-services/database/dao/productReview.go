package dao

import (
	"time"

	mp "github.com/geraldo-labs/merge-struct"
	gdatabase "github.com/prakashjegan/review-processor/app/database"
	"github.com/prakashjegan/review-processor/app/review-services/database/model"
	"github.com/prakashjegan/review-processor/app/review-services/utils"
	"gorm.io/gorm"
)

// TODO create Complete DAO for productReview.go files with all tables CRUD operations.
// Product Review Table Operations.
type DAO struct {
	DB *gorm.DB
}
type ProductReviewDAO struct {
	DAO
}

func NewProductReviewDAO(db *gorm.DB) *ProductReviewDAO {
	return &ProductReviewDAO{
		DAO: DAO{
			DB: db,
		},
	}
}

func (p *ProductReviewDAO) CreateOrUpdate(productReview *model.ProductReview) (*model.ProductReview, error) {
	db := gdatabase.GetDB()
	oldProductReview := &model.ProductReview{}
	err := db.Where("hotel_id = ?", productReview.HotelID).First(oldProductReview).Error
	if err != nil {
		oldProductReview.ID = utils.GetUID()
	}
	mp.Struct(oldProductReview, productReview)
	err = db.Save(oldProductReview).Error
	if err != nil {
		return nil, err
	}
	productReviewEvent := &model.ProductReviewEvent{}
	mp.Struct(productReviewEvent, oldProductReview)
	productReviewEvent.ID = utils.GetUID()
	productReviewEvent.ProductReviewId = oldProductReview.ID
	productReviewEvent.CreatedAt = time.Now()
	productReviewEvent.UpdatedAt = time.Now()
	err = db.Create(productReviewEvent).Error
	if err != nil {
		return nil, err
	}
	return oldProductReview, err
}

func (p *ProductReviewDAO) GetByProductId(productId string) ([]model.ProductReview, error) {

	db := gdatabase.GetDB()
	oldProductReview := []model.ProductReview{}
	err := db.Where("hotel_id = ? ", productId).Find(&oldProductReview).Error
	return oldProductReview, err
}

// Review Comment Table Operations.
type ReviewCommentDAO struct {
	DAO
}

func NewReviewCommentDAO(db *gorm.DB) *ReviewCommentDAO {
	return &ReviewCommentDAO{
		DAO: DAO{
			DB: db,
		},
	}
}

func (r *ReviewCommentDAO) CreateOrUpdate(reviewComment *model.ReviewComment) (*model.ReviewComment, error) {
	db := gdatabase.GetDB()
	oldReviewComment := &model.ReviewComment{}
	err := db.Where("hotel_review_id =? and hotel_id = ? ", reviewComment.HotelReviewID, reviewComment.HotelID).First(oldReviewComment).Error
	if err != nil {
		oldReviewComment.ID = utils.GetUID()
	}
	mp.Struct(oldReviewComment, reviewComment)
	err = db.Save(oldReviewComment).Error
	if err != nil {
		return nil, err
	}
	reviewCommentEvent := &model.ReviewCommentEvent{}
	mp.Struct(reviewCommentEvent, reviewComment)
	reviewCommentEvent.ID = utils.GetUID()
	reviewCommentEvent.ReviewCommentId = oldReviewComment.ID
	reviewCommentEvent.CreatedAt = time.Now()
	reviewCommentEvent.UpdatedAt = time.Now()
	err = db.Create(reviewCommentEvent).Error
	if err != nil {
		return nil, err
	}
	return oldReviewComment, err

}

func (r *ReviewCommentDAO) GetByReviewId(reviewId int64) ([]model.ReviewComment, error) {

	db := gdatabase.GetDB()
	oldReviewComment := []model.ReviewComment{}
	err := db.Where("product_review_id = ? ", reviewId).Find(&oldReviewComment).Error
	return oldReviewComment, err
}

// Reviewer Info Table Operations.
type ReviewerInfoDAO struct {
	DAO
}

func NewReviewerInfoDAO(db *gorm.DB) *ReviewerInfoDAO {
	return &ReviewerInfoDAO{
		DAO: DAO{
			DB: db,
		},
	}
}

func (r *ReviewerInfoDAO) CreateOrUpdate(reviewerInfo *model.ReviewerInfo) (*model.ReviewerInfo, error) {
	db := gdatabase.GetDB()
	oldReviewerInfo := &model.ReviewerInfo{}
	oldReviewerInfo.ID = utils.GetUID()
	err := db.Where("product_review_id =? and review_comment_id = ? ", reviewerInfo.ProductReviewId, reviewerInfo.ReviewCommentId).First(oldReviewerInfo).Error
	if err != nil {
		oldReviewerInfo.ID = utils.GetUID()
	}
	mp.Struct(oldReviewerInfo, reviewerInfo)
	err = db.Save(oldReviewerInfo).Error
	if err != nil {
		return nil, err
	}

	reviewerInfoEvent := &model.ReviewerInfoEvent{}
	mp.Struct(reviewerInfoEvent, reviewerInfo)
	reviewerInfoEvent.ID = utils.GetUID()
	reviewerInfoEvent.ReviewInfoId = oldReviewerInfo.ID
	reviewerInfoEvent.CreatedAt = time.Now()
	reviewerInfoEvent.UpdatedAt = time.Now()
	err = db.Create(reviewerInfoEvent).Error
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *ReviewerInfoDAO) GetByReviewCommentId(reviewCommentId string) ([]model.ReviewerInfo, error) {
	db := gdatabase.GetDB()
	oldReviewerInfo := []model.ReviewerInfo{}
	err := db.Where("review_comment_id = ? ", reviewCommentId).Find(&oldReviewerInfo).Error
	return oldReviewerInfo, err
}

// Review By Provider
type ProductReviewByProviderDAO struct {
	DAO
}

func NewProductReviewByProviderDAO(db *gorm.DB) *ProductReviewByProviderDAO {
	return &ProductReviewByProviderDAO{
		DAO: DAO{
			DB: db,
		},
	}
}

func (p *ProductReviewByProviderDAO) CreateOrUpdate(productReview *model.ProductReviewByProvider) (*model.ProductReviewByProvider, error) {
	db := gdatabase.GetDB()
	oldProductReview := &model.ProductReviewByProvider{}
	err := db.Where("hotel_id = ? and provider_id = ? ", productReview.HotelID, productReview.ProviderId).First(oldProductReview).Error
	if err != nil {
		oldProductReview.ID = utils.GetUID()
	}
	mp.Struct(oldProductReview, productReview)
	err = db.Save(oldProductReview).Error
	if err != nil {
		return nil, err
	}
	oldProductReviewByProviderEvent := &model.ProductReviewByProviderEvent{}
	mp.Struct(oldProductReviewByProviderEvent, oldProductReview)
	oldProductReviewByProviderEvent.ID = utils.GetUID()
	oldProductReviewByProviderEvent.ProductReviewByProviderId = oldProductReview.ID
	oldProductReviewByProviderEvent.CreatedAt = time.Now()
	oldProductReviewByProviderEvent.UpdatedAt = time.Now()
	err = db.Create(oldProductReviewByProviderEvent).Error
	if err != nil {
		return nil, err
	}
	return oldProductReview, err
}

func (p *ProductReviewByProviderDAO) GetByProductId(productId uint64) ([]model.ProductReviewByProvider, error) {
	db := gdatabase.GetDB()
	oldProductReview := []model.ProductReviewByProvider{}
	err := db.Where("product_review_id = ? ", productId).Find(&oldProductReview).Error
	return oldProductReview, err
}

type MetricData struct {
	OverAllPerformance float64 `json:"overallPerformance"`
	TotalReviewCount   int64   `json:"totalReviewCount"`
	ProviderId         int     `json:"providerId"`
	ProductReviewId    int     `json:"productReviewId"`
}

func (p *ProductReviewByProviderDAO) GetMetricsByProductId(productId uint64) ([]MetricData, error) {
	db := gdatabase.GetDB()
	metricData := make([]MetricData, 0)
	err := db.Model(&model.ProductReviewByProvider{}).Select(`avg(overall_score) as overall_performance,sum(review_count) as total_review_count, product_review_id`).Group("product_review_id").Where("product_review_id=?", productId).Scan(&metricData).Error
	return metricData, err
}

// product review  Grade table operations.
type ProductReviewGradeDAO struct {
	DAO
}

func NewProductReviewGradeDAO(db *gorm.DB) *ProductReviewGradeDAO {
	return &ProductReviewGradeDAO{
		DAO: DAO{
			DB: db,
		},
	}
}

func (p *ProductReviewGradeDAO) CreateOrUpdate(productReviewGrade *model.ProductReviewGrade) (*model.ProductReviewGrade, error) {
	db := gdatabase.GetDB()
	oldProductReviewGrade := &model.ProductReviewGrade{}
	err := db.Where("product_review_id = ? and provider_id = ? and grade_attribute = ? ", productReviewGrade.ProductReviewId, productReviewGrade.ProviderId, productReviewGrade.GradeAttribute).First(oldProductReviewGrade).Error
	if err != nil {
		oldProductReviewGrade.ID = utils.GetUID()
	}
	mp.Struct(oldProductReviewGrade, productReviewGrade)
	err = db.Save(oldProductReviewGrade).Error
	if err != nil {
		return nil, err
	}

	oldProductReviewGradeEvent := &model.ProductReviewGradeEvent{}
	mp.Struct(oldProductReviewGradeEvent, oldProductReviewGrade)
	oldProductReviewGradeEvent.ID = utils.GetUID()
	oldProductReviewGradeEvent.ProductReviewGradeId = oldProductReviewGrade.ID
	oldProductReviewGradeEvent.CreatedAt = time.Now()
	oldProductReviewGradeEvent.UpdatedAt = time.Now()
	err = db.Create(oldProductReviewGradeEvent).Error
	if err != nil {
		return nil, err
	}

	return oldProductReviewGrade, nil
}

func (p *ProductReviewGradeDAO) GetByProductId(productId uint64, providerId int) ([]*model.ProductReviewGrade, error) {

	db := gdatabase.GetDB()
	oldProductReviewGrade := []*model.ProductReviewGrade{}
	err := db.Where("product_review_id = ? and provider_id = ? ", productId, providerId).Find(&oldProductReviewGrade).Error
	return oldProductReviewGrade, err
}
