package gameusers

import (
	"github.com/guregu/null"
	"gitlab.com/kafa1942/take10dashboard/common"
	"time"
)

type GetAllRequest struct {
	MetaData []common.MetaData `json:"metaData,omitempty"`
}

type GetRequest struct {
	ID string `uri:"id" binding:"required,number"`
}

type UpdateRequest struct {
	UserID      int         `gorm:"column:user_id;primary_key" json:"user_id" binding:"required,number"`
	PhoneNumber null.String `gorm:"column:phoneNumber" json:"phoneNumber"`
	Name        null.String `gorm:"column:name" json:"name"`
	FbLastName  null.String `gorm:"column:fb_last_name" json:"fb_last_name"`
	FbFirstName null.String `gorm:"column:fb_first_name" json:"fb_first_name"`
	Email       string      `gorm:"column:email" json:"email"`
	UserStatus  null.String `gorm:"column:user_status" json:"user_status"`
}

type ProfileResponse struct {
	User
	UserProfile
	InvitePromoCount            int            `json:"invitePromoCount"`
	CompletePayment             string         `json:"completePayment"`
	LastPhoneVerificationDate   null.Time      `json:"lastPhoneVerificationDate"`
	LastPhoneVerificationStatus null.String    `json:"lastPhoneVerificationStatus"`
	LastLoggedInDate            time.Time      `json:"lastLoggedInDate"`
	LastActivities              []ActivityLog  `json:"lastActivities"`
	Winnings                    []RandomWinner `json:"winnings"`
}
