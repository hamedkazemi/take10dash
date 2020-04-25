package questions

import (
	"database/sql"
	"github.com/guregu/null"
	"gitlab.com/kafa1942/take10dashboard/common"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

// Models should only be concerned with database schema, more strict checking should be put in validator.
type Question struct {
	ID               int       `gorm:"column:id;primary_key" json:"id"`
	QuestionTypeID   string    `gorm:"column:questionType" json:"questionType"`
	QuestionType     QuestionsType `gorm:"foreignkey:QuestionID"`
	QuestionDiffType null.Int  `gorm:"column:questionDiffType" json:"questionDiffType"`
	CategoryID       int       `gorm:"column:category" json:"category"`
	Category 		 Category  `gorm:"foreignkey:ID;association_foreignkey:CategoryID"`
	Status           null.Int  `gorm:"column:status" json:"status"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdateAt         time.Time `gorm:"column:update_at" json:"update_at"`
	Answers          []Answers `gorm:"foreignkey:QuestionID;association_foreignkey:ID"`
}

// TableName sets the insert table name for this struct type
func (q *Question) TableName() string {
	return "Question"
}

type QuestionsType struct {
	ID         int         `gorm:"column:id;primary_key" json:"id"`
	QuestionID int         `gorm:"column:questionID" json:"questionID"`
	TextType   null.String `gorm:"column:textType" json:"textType"`
	ImageType  string      `gorm:"column:imageType" json:"imageType"`
}

// TableName sets the insert table name for this struct type
func (q *QuestionsType) TableName() string {
	return "QuestionsType"
}

type Category struct {
	ID   int    `gorm:"column:id;primary_key" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

// TableName sets the insert table name for this struct type
func (c *Category) TableName() string {
	return "Category"
}

type Answers struct {
	ID            int       `gorm:"column:id;primary_key" json:"id"`
	QuestionID    int       `gorm:"column:questionID" json:"questionID"`
	CorrectAnswer string    `gorm:"column:correctAnswer" json:"correctAnswer"`
	AnswerTypeString    string    `gorm:"column:answerType" json:"answerType"`
	AnswerType 	  AnswersType  `gorm:"foreignkey:AnswerID;association_foreignkey:AnswerTypeString"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdateAt      time.Time `gorm:"column:update_at" json:"update_at"`
}

// TableName sets the insert table name for this struct type
func (a *Answers) TableName() string {
	return "Answers"
}

type AnswersType struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	TextType  null.String `gorm:"column:textType" json:"textType"`
	ImageType null.String `gorm:"column:imageType" json:"imageType"`
	AnswerID  int         `gorm:"column:answerID" json:"answerID"`
}

// TableName sets the insert table name for this struct type
func (a *AnswersType) TableName() string {
	return "AnswersType"
}

// Migrate the schema of database if needed
func AutoMigrate() {
	db := common.GetDB()
	db.AutoMigrate(&Question{})
	db.AutoMigrate(&QuestionsType{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&Answers{})
	db.AutoMigrate(&AnswersType{})
}
