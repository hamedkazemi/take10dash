package questions

import (
	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

type QuestionSerializer struct {
	C *gin.Context
	Question
}

type QuestionsSerializer struct {
	C         *gin.Context
	Questions []Question
}

type AnswerSerializer struct {
	C *gin.Context
	Answers
}

type AnswersSerializer struct {
	C       *gin.Context
	Answers []Answers
}

type CategorySerializer struct {
	C *gin.Context
	Category
}

type CategoriesSerializer struct {
	C          *gin.Context
	Categories []Category
}

type AnswerResponse struct {
	ID            int         `json:"id"`
	QuestionId    int         `json:"questionId"`
	CorrectAnswer bool        `json:"correctAnswer"`
	AnswerType    string      `json:"answerType"`
	AnswerContext null.String `json:"answerContext"`
	CreatedAt     string      `json:"createdAt"`
	UpdatedAt     string      `json:"updatedAt"`
}

type QuestionResponse struct {
	ID                 int              `json:"id"`
	QuestionType       string           `json:"questionType"`    // image or text
	QuestionContext    null.String      `json:"questionContext"` // textType or imageType Value
	QuestionDifficulty null.Int         `json:"questionDifficulty"`
	CategoryId         int              `json:"categoryId"`
	CategoryTitle      string           `json:"categoryTitle"`
	Status             null.Int         `json:"status"`
	Answers            []AnswerResponse `json:"answers"`
	CreatedAt          string           `json:"createdAt"`
	UpdatedAt          string           `json:"updatedAt"`
}

type QuestionRequest struct {
	ID                 int             `json:"id,omitempty"`
	QuestionType       string          `json:"questionType"  binding:"required"`    // image or text
	QuestionContext    null.String     `json:"questionContext"  binding:"required"` // textType or imageType Value
	QuestionDifficulty null.Int        `json:"questionDifficulty"  binding:"required"`
	CategoryId         int             `json:"categoryId"  binding:"required"`
	Status             null.Int        `json:"status"  binding:"required"`
	Answers            []AnswerRequest `json:"answers"  binding:"required"`
}

type AnswerRequest struct {
	ID            int         `json:"id,omitempty"`
	CorrectAnswer bool        `json:"correctAnswer"  binding:"required"`
	AnswerType    string      `json:"answerType"  binding:"required"`
	AnswerContext null.String `json:"answerContext"  binding:"required"`
}

type QuestionsDeleteRequest struct {
	QuestionIds []int `json:"questionIds"  binding:"required"`
}

type QuestionGetRequest struct {
	ID string `uri:"id" binding:"required"`
}

type CategoriesResponse struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

func (s *AnswerSerializer) Response() AnswerResponse {
	isCorrect := false
	if s.CorrectAnswer == "true" {
		isCorrect = true
	}
	var answerContext null.String
	if s.AnswerTypeID == "text" {
		answerContext = s.AnswerType.TextType
	} else {
		answerContext = s.AnswerType.ImageType
	}
	response := AnswerResponse{
		ID:            s.ID,
		QuestionId:    s.QuestionID,
		CorrectAnswer: isCorrect,
		AnswerType:    s.AnswerTypeID,
		AnswerContext: answerContext,
		CreatedAt:     s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt:     s.UpdateAt.UTC().Format("2006-01-02T15:04:05.999Z"),
	}
	return response
}

func (s *AnswersSerializer) Response() []AnswerResponse {
	response := []AnswerResponse{}
	for _, answer := range s.Answers {
		serializer := AnswerSerializer{s.C, answer}
		response = append(response, serializer.Response())
	}
	return response
}

func (s *QuestionSerializer) Response() QuestionResponse {
	var questionContext null.String
	if s.QuestionTypeID == "text" {
		questionContext = s.QuestionType.TextType
	} else {
		questionContext = s.QuestionType.ImageType
	}
	AnswersSerializer := AnswersSerializer{s.C, s.Answers}
	response := QuestionResponse{
		ID:                 s.ID,
		QuestionType:       s.QuestionTypeID,
		QuestionContext:    questionContext,
		QuestionDifficulty: s.QuestionDiffType,
		CategoryId:         s.CategoryID,
		CategoryTitle:      s.Category.Name,
		Status:             s.Status,
		Answers:            AnswersSerializer.Response(),
		CreatedAt:          s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt:          s.UpdateAt.UTC().Format("2006-01-02T15:04:05.999Z"),
	}
	return response
}

func (s *QuestionsSerializer) Response() []QuestionResponse {
	response := []QuestionResponse{}
	for _, question := range s.Questions {
		serializer := QuestionSerializer{s.C, question}
		response = append(response, serializer.Response())
	}
	return response
}

func (s *CategorySerializer) Response() CategoriesResponse {
	response := CategoriesResponse{
		Label: s.Name,
		Value: s.ID,
	}
	return response
}

func (s *CategoriesSerializer) Response() []CategoriesResponse {
	response := []CategoriesResponse{}
	for _, category := range s.Categories {
		serializer := CategorySerializer{s.C, category}
		response = append(response, serializer.Response())
	}
	return response
}
