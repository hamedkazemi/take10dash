package questions

import (
	"database/sql"
	"github.com/guregu/null"
	"github.com/jinzhu/gorm"
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
	ID               int           `gorm:"column:id;primary_key" json:"id"`
	QuestionTypeID   string        `gorm:"column:questionType" json:"questionType"`
	QuestionType     QuestionsType `gorm:"foreignkey:ID;association_foreignkey:QuestionID"`
	QuestionDiffType null.Int      `gorm:"column:questionDiffType" json:"questionDiffType"`
	CategoryID       int           `gorm:"column:category" json:"category"`
	Category         Category      `gorm:"foreignkey:CategoryID;association_foreignkey:ID"`
	Status           null.Int      `gorm:"column:status" json:"status"`
	CreatedAt        time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdateAt         time.Time     `gorm:"column:update_at" json:"update_at"`
	Answers          []Answers     `gorm:"foreignkey:QuestionID;association_foreignkey:ID"`
}

// TableName sets the insert table name for this struct type
func (q *Question) TableName() string {
	return "Question"
}

// AfterDelete hook defined for cascade delete
func (question *Question) AfterDelete(tx *gorm.DB) error {
	common.PrettyPrint(question)
	res := tx.Exec("DELETE Answers,AnswersType FROM Answers "+
		"INNER JOIN AnswersType WHERE "+
		"AnswersType.answerID = Answers.id AND Answers.questionID = ? ", question.ID)
	if len(res.GetErrors()) > 0 {
		return res.Error
	}
	qt := QuestionsType{}
	res = tx.Unscoped().Where(QuestionsType{QuestionID: question.ID}).Delete(&qt)
	if len(res.GetErrors()) > 0 {
		return res.Error
	}
	return nil
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
	ID            int         `gorm:"column:id;primary_key" json:"id"`
	QuestionID    int         `gorm:"column:questionID" json:"questionID"`
	CorrectAnswer string      `gorm:"column:correctAnswer" json:"correctAnswer"`
	AnswerTypeID  string      `gorm:"column:answerType"`
	AnswerType    AnswersType `gorm:"foreignkey:ID;association_foreignkey:AnswerID;PRELOAD:true"`
	CreatedAt     time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdateAt      time.Time   `gorm:"column:update_at" json:"update_at"`
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
//func AutoMigrate() {
//	db := common.GetDB()
//	db.AutoMigrate(&Question{})
//	db.AutoMigrate(&QuestionsType{})
//	db.AutoMigrate(&Category{})
//	db.AutoMigrate(&Answers{})
//	db.AutoMigrate(&AnswersType{})
//}

func FindManyQuestions(filters []common.Filter, orderBy string, orderType string, limit int, offset int) ([]Question, int, error) {
	db := common.GetDB()
	var models []Question
	var count int

	var searchKeyword string
	//var conditions []string
	var whereCategory string
	var whereType string
	var whereDiff string

	for _, f := range filters {
		if f.Field == "textType" {
			searchKeyword = f.Value
		}
		if f.Field == "category" {
			whereCategory = "category = " + f.Value + " "
		}
		if f.Field == "questionType" {
			whereType = "questionType = '" + f.Value + "' "
		}
		if f.Field == "questionDiffType" {
			whereDiff = "questionDiffType = '" + f.Value + "' "
		}
	}

	var err error

	dbq := db.Table("Question").Select("DISTINCT Question.*").
		Joins("left join QuestionsType on Question.id = QuestionsType.questionID").
		Where(whereCategory).
		Where(whereType).
		Where(whereDiff).
		Joins("left join Answers on Question.id = Answers.questionID").
		Joins("left join AnswersType on Answers.id = AnswersType.answerID").
		Where("QuestionsType.textType LIKE '%"+searchKeyword+"%' "+
			"OR QuestionsType.imageType LIKE '%"+searchKeyword+"%' "+
			"OR AnswersType.textType LIKE '%"+searchKeyword+"%' "+
			"OR AnswersType.imageType LIKE '%"+searchKeyword+"%'").
		Set("gorm:auto_preload", true)

	err = dbq.Select("count(DISTINCT(Question.id))").Count(&count).Error
	err = dbq.Order(orderBy + " " + orderType).
		Limit(limit).
		Offset(offset).Find(&models).Error

	return models, count, err
}

func FindManyCategories() ([]Category, error) {
	db := common.GetDB()
	var models []Category
	err := db.Find(&models).Error
	return models, err
}
