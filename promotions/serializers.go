package promotions

import (
	"github.com/guregu/null"
	"gitlab.com/kafa1942/take10dashboard/common"
)

type GetAllRequest struct {
	MetaData []common.MetaData `json:"metaData,omitempty"`
}

type GetRequest struct {
	ID string `uri:"id" binding:"required,number"`
}

type PromotionCreateRequest struct {
	Code          string      `gorm:"column:code" json:"code" searchable:"string" binding:"required"`
	Points        null.Float  `gorm:"column:points" json:"points" searchable:"int"`
	StartDate     null.Time   `gorm:"column:start_date" json:"start_date" searchable:"false" validate:"required,datetime"`
	EndDate       null.Time   `gorm:"column:end_date" json:"end_date" searchable:"false" validate:"required,datetime"`
	PromoCategory null.String `gorm:"column:promo_category" json:"promo_category"  searchable:"string"`
	PrizeURL      null.String `gorm:"column:prize_url" json:"prize_url"  searchable:"string" validate:"url"`
	Status        string      `gorm:"column:status" json:"status"  searchable:"string" binding:"required"`
}

type PromotionUpdateRequest struct {
	ID            int         `gorm:"column:id;primary_key" json:"id" searchable:"int" binding:"required"`
	Code          string      `gorm:"column:code" json:"code" searchable:"string" binding:"required"`
	Points        null.Float  `gorm:"column:points" json:"points" searchable:"int"`
	StartDate     null.Time   `gorm:"column:start_date" json:"start_date" searchable:"false" validate:"required,datetime"`
	EndDate       null.Time   `gorm:"column:end_date" json:"end_date" searchable:"false" validate:"required,datetime"`
	PromoCategory null.String `gorm:"column:promo_category" json:"promo_category"  searchable:"string"`
	PrizeURL      null.String `gorm:"column:prize_url" json:"prize_url"  searchable:"string" validate:"url"`
	Status        string      `gorm:"column:status" json:"status"  searchable:"string"  binding:"required"`
}
