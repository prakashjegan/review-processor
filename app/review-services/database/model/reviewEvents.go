package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductReviewEvent struct {
	ID              uint64 `gorm:"primaryKey" json:"id,omitempty"`
	ProductReviewId uint64 `gorm:"index" json:"productReviewId,omitempty"`
	HotelID         int    `json:"hotelId,omitempty"`
	HotelName       string `json:"hotelName,omitempty"`

	OverallScore int64 `json:"overallScore,omitempty"` //( average of all provider score . score Multiplied By 100 to make it uint64 )
	ReviewCount  int64 `json:"reviewCount,omitempty"`  // total of all provider review count.

	LastReviewDate time.Time `json:"reviewDate,omitempty"`
	LastComment    string    `json:"lastComment,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted bool           `gorm:"-" json:"isDeleted,omitempty"`
	DeletedBy string         `gorm:"-" json:"deletedBy,omitempty"`
}

type ReviewCommentEvent struct {
	ID                   uint64 `gorm:"primaryKey" json:"id,omitempty"`
	ReviewCommentId      uint64 `gorm:"index" json:"reviewCommentId,omitempty"`
	ProductReviewId      uint64 `gorm:"index" json:"productReviewId,omitempty"`
	ProductReviewEventId uint64 `gorm:"index" json:"productReviewEventId,omitempty"`

	HotelID   int    `json:"hotelId,omitempty"`
	HotelName string `json:"hotelName,omitempty"`

	IsShowReviewResponse  bool    `json:"isShowReviewResponse"`
	HotelReviewID         int     `gorm:"index" json:"hotelReviewId"`
	ProviderID            int     `gorm:"index" json:"providerId"`
	Rating                float64 `json:"rating"`
	CheckInDateMonthYear  string  `json:"checkInDateMonthAndYear"`
	EncryptedReviewData   string  `json:"encryptedReviewData"`
	FormattedRating       string  `json:"formattedRating"`
	FormattedReviewDate   string  `json:"formattedReviewDate"`
	RatingText            string  `json:"ratingText"`
	ResponderName         string  `json:"responderName"`
	ResponseDateText      string  `json:"responseDateText"`
	ResponseTranslateSrc  string  `json:"responseTranslateSource"`
	ReviewComments        string  `json:"reviewComments"`
	ReviewNegatives       string  `json:"reviewNegatives"`
	ReviewPositives       string  `json:"reviewPositives"`
	ReviewProviderLogo    string  `json:"reviewProviderLogo"`
	ReviewProviderText    string  `json:"reviewProviderText"`
	ReviewTitle           string  `json:"reviewTitle"`
	TranslateSource       string  `json:"translateSource"`
	TranslateTarget       string  `json:"translateTarget"`
	ReviewDate            string  `json:"reviewDate"` // Use time.Time if parsing
	OriginalTitle         string  `json:"originalTitle"`
	OriginalComment       string  `json:"originalComment"`
	FormattedResponseDate string  `json:"formattedResponseDate"`

	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted bool           `gorm:"-" json:"isDeleted,omitempty"`
	DeletedBy string         `gorm:"-" json:"deletedBy,omitempty"`
}

type ReviewerInfoEvent struct {
	ID                   uint64 `gorm:"primaryKey" json:"id,omitempty"`
	ReviewInfoId         uint64 `gorm:"index" json:"reviewInfoId,omitempty"`
	ReviewCommentId      uint64 `gorm:"index" json:"reviewCommentId,omitempty"`
	ReviewCommentEventId uint64 `gorm:"index" json:"reviewCommentEventId,omitempty"`
	ProductReviewEventId uint64 `gorm:"index" json:"productReviewEventId,omitempty"`
	ProductReviewId      uint64 `gorm:"index" json:"productReviewId,omitempty"`

	CountryName           string `json:"countryName"`
	DisplayMemberName     string `json:"displayMemberName"`
	FlagName              string `json:"flagName"`
	ReviewGroupName       string `json:"reviewGroupName"`
	RoomTypeName          string `json:"roomTypeName"`
	CountryID             int    `json:"countryId"`
	LengthOfStay          int    `json:"lengthOfStay"`
	ReviewGroupID         int    `json:"reviewGroupId"`
	RoomTypeID            int    `json:"roomTypeId"`
	ReviewerReviewedCount int    `json:"reviewerReviewedCount"`
	IsExpertReviewer      bool   `json:"isExpertReviewer"`
	IsShowGlobalIcon      bool   `json:"isShowGlobalIcon"`
	IsShowReviewedCount   bool   `json:"isShowReviewedCount"`

	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted bool           `gorm:"-" json:"isDeleted,omitempty"`
	DeletedBy string         `gorm:"-" json:"deletedBy,omitempty"`
}

type ProductReviewByProviderEvent struct {
	ID                        uint64 `gorm:"primaryKey" json:"id,omitempty"`
	ProductReviewByProviderId uint64 `gorm:"index" json:"productReviewByProviderId,omitempty"`
	ProductReviewId           uint64 `gorm:"index" json:"productReviewId,omitempty"`
	ProductReviewEventId      uint64 `gorm:"index" json:"productReviewEventId,omitempty"`

	HotelID   int    `json:"hotelId,omitempty"`
	HotelName string `json:"hotelName,omitempty"`

	Platform             string `gorm:"index" json:"platform,omitempty"`
	ProviderId           int    `gorm:"index" json:"providerId,omitempty"`
	ThirdPartyProviderId uint64 `gorm:"index" json:"thirdPartyProviderId,omitempty"`
	Provider             string `gorm:"index" json:"provider,omitempty"`
	OverallScore         int64  `json:"overallScore,omitempty"` //( score Multiplied By 100 to make it uint64 )
	ReviewCount          int64  `json:"reviewCount,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted bool           `gorm:"-" json:"isDeleted,omitempty"`
	DeletedBy string         `gorm:"-" json:"deletedBy,omitempty"`
}

type ProductReviewGradeEvent struct {
	ID uint64 `gorm:"primaryKey" json:"id,omitempty"`

	ProductReviewGradeId           uint64 `gorm:"index" json:"productReviewGradeId,omitempty"`
	ProductReviewByProviderId      uint64 `gorm:"index" json:"productReviewByProviderId,omitempty"`
	ProductReviewByProviderEventId uint64 `gorm:"index" json:"productReviewByProviderEventId,omitempty"`
	ProductReviewEventId           uint64 `gorm:"index" json:"productReviewEventId,omitempty"`
	ProductReviewId                uint64 `gorm:"index" json:"productReviewId,omitempty"`

	HotelID              int    `gorm:"index" json:"hotelId,omitempty"`
	HotelName            string `gorm:"index" json:"hotelName,omitempty"`
	Platform             string `gorm:"index" json:"platform,omitempty"`
	ProviderId           int    `gorm:"index" json:"providerId,omitempty"`
	ThirdPartyProviderId uint64 `gorm:"index" json:"thirdPartyProviderId,omitempty"`
	Provider             string `gorm:"index" json:"provider,omitempty"`

	GradeAttribute string `json:"gradeAttribute,omitempty"`
	Grade          int64  `json:"grade,omitempty"` // out of 10 multiple by 100 to make it integer.

	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted bool           `gorm:"-" json:"isDeleted,omitempty"`
	DeletedBy string         `gorm:"-" json:"deletedBy,omitempty"`
}
