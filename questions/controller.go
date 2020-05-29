package questions

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"gitlab.com/kafa1942/take10dashboard/common"
	_ "gitlab.com/kafa1942/take10dashboard/users"
	"net/http"
	"strconv"
	"time"
)

func GetAllQuestions(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	searchQuery := c.Query("query")
	categoryID := c.Query("category")
	questionType := c.Query("questionType")
	difficultyLevel := c.Query("difficultyLevel")
	orderBy := c.Query("orderBy")
	orderType := c.Query("orderType")

	// default values for limit and offset
	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}
	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		logrus.Error("Question", err)
		c.JSON(http.StatusNotFound, common.NewError("questions", errors.New("Invalid param")))
		return
	}

	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		logrus.Error("Question", err)
		c.JSON(http.StatusNotFound, common.NewError("questions", errors.New("Invalid param")))
		return
	}

	var filters []common.Filter

	if orderBy == "" {
		orderBy = "created_at"
	}

	if orderType == "" {
		orderType = "DESC"
	}

	if searchQuery != "" {
		filters = append(filters, common.Filter{Field: "textType", Operator: "LIKE", Value: searchQuery})
	}

	if categoryID != "" {
		filters = append(filters, common.Filter{Field: "category", Operator: "=", Value: categoryID})
	}

	if questionType != "" {
		filters = append(filters, common.Filter{Field: "questionType", Operator: "=", Value: questionType})
	}

	if difficultyLevel != "" {
		filters = append(filters, common.Filter{Field: "questionDiffType", Operator: "=", Value: difficultyLevel})
	}

	questionsModel, modelCount, err := FindManyQuestions(filters, orderBy, orderType, limit_int, offset_int)
	if err != nil {
		logrus.Error("Question", err)
		c.JSON(http.StatusNotFound, common.NewError("questions", errors.New("Something went wrong!")))
		return
	}

	serializer := QuestionsSerializer{c, questionsModel}

	var metaData common.MetaData
	metaData.Pagination.Count = modelCount
	metaData.Pagination.Limit = limit_int
	metaData.Pagination.Offset = offset_int
	metaData.Filters = filters
	metaData.Order.OrderBy = orderBy
	metaData.Order.OrderType = orderType

	c.JSON(http.StatusOK, gin.H{"data": serializer.Response(), "metaData": metaData})
}

func GetQuestion(c *gin.Context) {
	var qr QuestionGetRequest
	if err := c.ShouldBindUri(&qr); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("message", errors.New(err.Error())))
		return
	}
	var question Question
	common.GetDB().Set("gorm:auto_preload", true).First(&question, qr.ID)
	serializer := QuestionSerializer{c, question}
	c.JSON(http.StatusOK, gin.H{"data": serializer.Response()})
}

func GetAllCategories(c *gin.Context) {
	categories, err := FindManyCategories()
	if err != nil {
		logrus.Error("Category", err)
		c.JSON(http.StatusNotFound, common.NewError("Category", errors.New("Something went wrong!")))
		return
	}
	serializer := CategoriesSerializer{c, categories}
	c.JSON(http.StatusOK, gin.H{"data": serializer.Response()})
}

func UpdateQuestion(c *gin.Context) {
	var input QuestionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("Update Question", err))
		return
	}
	// Check if question id is provided
	if input.ID == 0 {
		c.JSON(http.StatusBadRequest, common.NewError("Update Question", errors.New("ID is required for update.")))
		return
	}
	// Check if question exists and load model with data
	var question Question
	if common.GetDB().Set("gorm:auto_preload", true).First(&question, input.ID).RecordNotFound() {
		c.JSON(http.StatusBadRequest, common.NewError("Update Question", errors.New("Question id ["+fmt.Sprint(input.ID)+"] not found.")))
		return
	}

	ct := 0
	// validation for Answer Ids
	for _, v := range input.Answers {
		if v.ID == 0 {
			c.JSON(http.StatusBadRequest, common.NewError("Update Question", errors.New("Answers IDs are required for update.")))
			return
		}
		if v.CorrectAnswer == true {
			ct++
		}
	}

	// One true answer is acceptable
	if ct > 1 {
		c.JSON(http.StatusBadRequest, common.NewError("Update Question", errors.New("Only one true answer is acceptable.")))
		return
	}

	if len(input.Answers) < 4 {
		c.JSON(http.StatusBadRequest, common.NewError("Update Question", errors.New("Invalid answers input, expected 4 items.")))
		return
	}

	var newCategory Category
	// get new category if needed
	cnf := common.GetDB().First(&newCategory, input.CategoryId).RecordNotFound()
	if question.CategoryID != input.CategoryId {
		if cnf {
			c.JSON(http.StatusBadRequest, common.NewError("Update Question", errors.New("Selected category not found")))
			return
		}
	}

	// Now it's time to update the question
	question.CategoryID = input.CategoryId
	question.Category.ID = input.CategoryId
	question.Category.Name = newCategory.Name
	question.QuestionDiffType = input.QuestionDifficulty
	question.QuestionTypeID = input.QuestionType
	question.Status = input.Status

	question.UpdateAt = time.Now()

	if question.QuestionTypeID == "text" {
		question.QuestionType.QuestionID = input.ID
		question.QuestionType.TextType = input.QuestionContext
		question.QuestionType.ImageType = ""
	} else {
		question.QuestionType.QuestionID = input.ID
		if input.QuestionContext.Valid { // bad database design
			question.QuestionType.ImageType = input.QuestionContext.String
		}
		question.QuestionType.TextType = null.String{}
	}

	// get answers
	var answers []Answers
	common.GetDB().Set("gorm:auto_preload", true).Where(Answers{QuestionID: input.ID}).Find(&answers)
	for i, v := range answers {
		if v.ID == input.Answers[i].ID {
			v.ID = input.Answers[i].ID
			v.QuestionID = input.ID
			v.UpdateAt = time.Now()
			v.AnswerTypeID = input.Answers[i].AnswerType
			v.CorrectAnswer = fmt.Sprint(input.Answers[i].CorrectAnswer)
			v.AnswerType.AnswerID = input.Answers[i].ID
			// update AnswerType
			if input.Answers[i].AnswerType == "text" {
				v.AnswerType.TextType = input.Answers[i].AnswerContext
				v.AnswerType.ImageType = null.String{}
			} else {
				v.AnswerType.ImageType = input.Answers[i].AnswerContext
				v.AnswerType.TextType = null.String{}
			}
			common.GetDB().Save(&v)
			question.Answers[i] = v
		}

	}

	err := common.GetDB().Set("gorm:association_autoupdate", true).Save(&question).Error
	if err != nil {
		logrus.Error("Update Question", err)
		c.JSON(http.StatusNotFound, common.NewError("Update Question", errors.New("Something went wrong!")))
		return
	}
	serializer := QuestionSerializer{c, question}
	c.JSON(http.StatusOK, gin.H{"data": serializer.Response(), "message": "Question Updated Successfully"})
}

func CreateQuestion(c *gin.Context) {
	var input QuestionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("Update Question", err))
		return
	}
}

func DeleteQuestion(c *gin.Context) {
	var input QuestionsDeleteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("message", errors.New("Invalid Input.")))
		return
	}
	for _, v := range input.QuestionIds {
		if v != 0 {
			question := Question{}
			question.ID = v
			res := common.GetDB().Unscoped().Delete(&question)
			if len(res.GetErrors()) > 0 {
				logrus.Error("DeleteQuestion", res.GetErrors(), res.Error)
				c.JSON(http.StatusOK, gin.H{"message": "Can't Delete the Question id " + fmt.Sprint(v) + ", see the logs for more information."})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Question Deleted Successfully"})
}

func GetStatistics(c *gin.Context) {
	// Raw SQL
	var StatResponse []QuestionStatsResponse
	rows, err := common.GetDB().Select("category.name as category,(SELECT COUNT(*) FROM `Question` WHERE `questionType` = 'text' AND category = category.id) as textcount,(SELECT COUNT(*) FROM `Question` WHERE `questionType` = 'picture' AND category = category.id) as picturecount,(SELECT COUNT(*) FROM `Question` WHERE `questionDiffType` = 1 AND category = category.id) as level1,(SELECT COUNT(*) FROM `Question` WHERE `questionDiffType` = 2 AND category = category.id) as level2,(SELECT COUNT(*) FROM `Question` WHERE `questionDiffType` = 3 AND category = category.id) as level3,(SELECT COUNT(*) FROM `Question` WHERE `questionDiffType` = 4 AND category = category.id) as level4,(SELECT COUNT(*) FROM `Question` WHERE `questionDiffType` = 5 AND category = category.id) as level5").Table("category").Rows() // (*s
	if err != nil {
		logrus.Error("Question Stat", err.Error)
		c.JSON(http.StatusOK, gin.H{"message": "Can't get stats, see the logs for more information."})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var stat QuestionStatsResponse
		err := common.GetDB().ScanRows(rows, &stat)
		fmt.Println(stat)
		StatResponse = append(StatResponse, stat)
		if err != nil {
			logrus.Error("Question Stat", err.Error)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": StatResponse})

}
