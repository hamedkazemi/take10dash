package questions

import (
	"github.com/gin-gonic/gin"
)

func ConfigGinRouter(router gin.IRoutes) {
	configGinAnswersRouter(router)
	configGinAnswersTypesRouter(router)
	configGinCategoriesRouter(router)
	configGinQuestionsRouter(router)
	configGinQuestionsTypesRouter(router)
	return
}

func configGinAnswersRouter(router gin.IRoutes) {
	router.GET("/answers", GetAllAnswers)
	router.POST("/answers", AddAnswer)
	router.GET("/answers/:id", GetAnswer)
	router.PUT("/answers/:id", UpdateAnswer)
	router.DELETE("/answers/:id", DeleteAnswer)
}

func configGinAnswersTypesRouter(router gin.IRoutes) {
	//
}

func configGinCategoriesRouter(router gin.IRoutes) {
	//
}

func configGinQuestionsRouter(router gin.IRoutes) {
	//
}

func configGinQuestionsTypesRouter(router gin.IRoutes) {
	//
}
