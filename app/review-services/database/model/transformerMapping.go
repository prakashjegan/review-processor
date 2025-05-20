package model

type TransformerMapping struct {
	ID                 uint64 `gorm:"primaryKey" json:"id,omitempty"`
	TransformerName    string `json:"transformerName,omitempty"`
	TransformedColumns string `json:"transformedColumns,omitempty"`

	ThirdPartyName             string `json:"thirdPartyName,omitempty"`
	ThirdPartyConfigType       string `json:"thirdPartyConfigType,omitempty"`
	ThirdPartyConnectionConfig string `json:"thirdPartyConnectionConfig,omitempty"` // Json Structure
	ThirdPartyReviewConfig     string `json:"thirdPartyReviewConfig,omitempty"`     // Json Structure

}
