package promotions

import (
	"database/sql"
	"github.com/guregu/null"
	_ "github.com/jinzhu/gorm"
	"gitlab.com/kafa1942/take10dashboard/common"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

// Models should only be concerned with database schema, more strict checking should be put in validator.

type Promotion struct {
	ID            int         `gorm:"column:id;primary_key" json:"id" searchable:"int"`
	Code          string      `gorm:"column:code" json:"code" searchable:"string"`
	Points        null.Float  `gorm:"column:points" json:"points" searchable:"int"`
	StartDate     null.Time   `gorm:"column:start_date" json:"start_date" searchable:"false"`
	EndDate       null.Time   `gorm:"column:end_date" json:"end_date" searchable:"false"`
	PromoCategory null.String `gorm:"column:promo_category" json:"promo_category"  searchable:"string"`
	PrizeURL      null.String `gorm:"column:prize_url" json:"prize_url"  searchable:"string"`
	Status        string      `gorm:"column:status" json:"status"  searchable:"string"`
}

// TableName sets the insert table name for this struct type
func (p *Promotion) TableName() string {
	return "Promotion"
}

func FindMany(filters []common.Filter, orderBy string, orderType string, limit int, offset int) ([]Promotion, int, error) {
	var models []Promotion
	var count int

	db := common.GetDB().Model(&models)

	dbq := common.BuildSearchByTags(Promotion{}, filters, db)

	var err error

	err = dbq.Count(&count).Error
	err = dbq.Order(orderBy + " " + orderType).
		Limit(limit).
		Offset(offset).Find(&models).Error

	return models, count, err
}
