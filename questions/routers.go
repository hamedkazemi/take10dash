package questions

import (
	"github.com/gin-gonic/gin"
)

func ConfigGinRouter(router gin.IRoutes) {
	configGinCategoriesRouter(router)
	configGinQuestionsRouter(router)
	return
}

func configGinCategoriesRouter(router gin.IRoutes) {
	router.GET("/categories", GetAllCategories)
}

func configGinQuestionsRouter(router gin.IRoutes) {
	router.GET("/questions", GetAllQuestions)
	router.GET("/questions/:id", GetQuestion)
	router.PATCH("/questions", UpdateQuestion)
	router.POST("/questions", CreateQuestion)
	router.DELETE("/questions", DeleteQuestion)

}
