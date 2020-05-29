package questions

import (
	"github.com/gin-gonic/gin"
)

func ConfigGinRouter(router *gin.RouterGroup) {
	// categories route
	configGinCategoriesRouter(router)

	// question module group
	r := router.Group("/questions")
	{
		configGinQuestionsRouter(r)
	}

}

func configGinCategoriesRouter(router gin.IRoutes) {
	router.GET("/categories/", GetAllCategories)
	router.GET("/stats/questions/", GetStatistics)
}

func configGinQuestionsRouter(router gin.IRoutes) {
	// routes
	router.GET("/", GetAllQuestions)
	router.GET("/:id", GetQuestion)
	router.PATCH("/", UpdateQuestion)
	router.POST("/", CreateQuestion)
	router.DELETE("/", DeleteQuestion)
}
