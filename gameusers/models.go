package gameusers

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

type User struct {
	UserID        int         `gorm:"column:user_id;primary_key" json:"user_id"searchable:"int"`
	IsVerified    null.Int    `gorm:"column:isVerified" json:"isVerified"searchable:"int"`
	PhoneNumber   null.String `gorm:"column:phoneNumber" json:"phoneNumber"searchable:"string"`
	Name          null.String `gorm:"column:name" json:"name"searchable:"string"`
	Gender        null.String `gorm:"column:gender" json:"gender"searchable:"string"`
	SocialAccount null.String `gorm:"column:socialAccount" json:"socialAccount"searchable:"string"`
	FbLastName    null.String `gorm:"column:fb_last_name" json:"fb_last_name"searchable:"string"`
	FacebookID    null.String `gorm:"column:facebookId" json:"facebookId"searchable:"string"`
	FbFirstName   null.String `gorm:"column:fb_first_name" json:"fb_first_name"searchable:"string"`
	Email         string      `gorm:"column:email" json:"email"searchable:"string"`
	UserStatus    null.String `gorm:"column:user_status" json:"user_status"searchable:"string"`
	UserRole      int         `gorm:"column:user_role" json:"user_role"searchable:"int"`
	CreatedAt     null.Time   `gorm:"column:created_at" json:"created_at"searchable:"datetime"`
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "Users"
}

type UserProfile struct {
	UserID                  int         `gorm:"column:user_id;primary_key" json:"user_id"`
	NoOfGamesPlayed         null.Int    `gorm:"column:no_of_games_played" json:"no_of_games_played"`
	NumOfCorrectAnswers     null.Int    `gorm:"column:num_of_correct_answers" json:"num_of_correct_answers"`
	NoOfPGamesPlayed        null.Int    `gorm:"column:no_of_p_games_played" json:"no_of_p_games_played"`
	NumOfTry                null.Int    `gorm:"column:num_of_try" json:"num_of_try"`
	NumOfGameTried          null.Int    `gorm:"column:num_of_game_tried" json:"num_of_game_tried"`
	NoOfPWinning            null.Int    `gorm:"column:no_of_p_winning" json:"no_of_p_winning"`
	NumOfWrongAnswers       null.Int    `gorm:"column:num_of_wrong_answers" json:"num_of_wrong_answers"`
	NoOfWinning             null.Int    `gorm:"column:no_of_winning" json:"no_of_winning"`
	GiftcardAccepted        null.Int    `gorm:"column:giftcard_accepted" json:"giftcard_accepted"`
	RandomPrize             int         `gorm:"column:random_prize" json:"random_prize"`
	FacebookToken           null.String `gorm:"column:facebook_token" json:"facebook_token"`
	Points                  null.Float  `gorm:"column:points" json:"points"`
	ActivePromocode         null.String `gorm:"column:active_promocode" json:"active_promocode"`
	ActivePromocodeCategory null.String `gorm:"column:active_promocode_category" json:"active_promocode_category"`
	IsBlocked               null.Int    `gorm:"column:isBlocked" json:"isBlocked"`
	FeedbackCredit          null.Int    `gorm:"column:feedbackCredit" json:"feedbackCredit"`
	FbShareCredit           null.Int    `gorm:"column:fbShareCredit" json:"fbShareCredit"`
	BlockedTime             null.Time   `gorm:"column:BlockedTime" json:"BlockedTime"`
	GameStatus              int         `gorm:"column:game_status" json:"game_status"`
	QuestionDifficultyLevel null.Int    `gorm:"column:question_difficulty_level" json:"question_difficulty_level"`
}

// TableName sets the insert table name for this struct type
func (u *UserProfile) TableName() string {
	return "UserProfile"
}

type PhoneVerification struct {
	ID               int         `gorm:"column:id;primary_key" json:"id"`
	PhoneNumber      string      `gorm:"column:phoneNumber" json:"phoneNumber"`
	DeviceUID        string      `gorm:"column:device_uid" json:"device_uid"`
	RequestID        null.String `gorm:"column:request_id" json:"request_id"`
	VerificationCode null.String `gorm:"column:verificationCode" json:"verificationCode"`
	IPAddress        null.String `gorm:"column:ip_address" json:"ip_address"`
	Carrier          null.String `gorm:"column:carrier" json:"carrier"`
	UserID           null.Int    `gorm:"column:user_id" json:"user_id"`
	Status           null.String `gorm:"column:status" json:"status"`
	CreatedAt        null.Time   `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (p *PhoneVerification) TableName() string {
	return "PhoneVerification"
}

type CompletePayment struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`
	UserID     int       `gorm:"column:userId" json:"userId"`
	Token      string    `gorm:"column:token" json:"token"`
	QuestionID int       `gorm:"column:questionId" json:"questionId"`
	Amount     string    `gorm:"column:amount" json:"amount"`
	CreatedAt  null.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (c *CompletePayment) TableName() string {
	return "CompletePayment"
}

type RandomWinner struct {
	ID          int         `gorm:"column:id;primary_key" json:"id"`
	UserID      int         `gorm:"column:user_id" json:"user_id"`
	Email       string      `gorm:"column:email" json:"email"`
	Date        null.String `gorm:"column:date" json:"date"`
	PhoneNumber null.String `gorm:"column:phoneNumber" json:"phoneNumber"`
	Category    null.String `gorm:"column:category" json:"category"`
	PrizeURL    string      `gorm:"column:prizeUrl" json:"prizeUrl"`
}

// TableName sets the insert table name for this struct type
func (r *RandomWinner) TableName() string {
	return "RandomWinners"
}

type InvitePromo struct {
	ID             int         `gorm:"column:id;primary_key" json:"id"`
	SenderUserID   int         `gorm:"column:sender_user_id" json:"sender_user_id"`
	SenderIP       null.String `gorm:"column:sender_ip" json:"sender_ip"`
	ReceiverUserID null.Int    `gorm:"column:receiver_user_id" json:"receiver_user_id"`
	ReceiverIP     null.String `gorm:"column:receiver_ip" json:"receiver_ip"`
	Promocode      string      `gorm:"column:promocode" json:"promocode"`
	Status         int         `gorm:"column:status" json:"status"`
	CreatedAt      time.Time   `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (i *InvitePromo) TableName() string {
	return "InvitePromo"
}

type ActivityLog struct {
	ID              int64       `gorm:"column:id;primary_key" json:"id"`
	UserID          int         `gorm:"column:user_id" json:"user_id"`
	IPAddress       null.String `gorm:"column:ip_address" json:"ip_address"`
	ActivityMessage string      `gorm:"column:activity_message" json:"activity_message"`
	CreatedAt       time.Time   `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (a *ActivityLog) TableName() string {
	return "ActivityLog"
}

func FindMany(filters []common.Filter, orderBy string, orderType string, limit int, offset int) ([]User, int, error) {
	var models []User
	var count int

	db := common.GetDB().Model(&models)

	dbq := common.BuildSearchByTags(User{}, filters, db)

	var err error

	err = dbq.Count(&count).Error
	err = dbq.Order(orderBy + " " + orderType).
		Limit(limit).
		Offset(offset).Find(&models).Error

	return models, count, err
}
