package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"gitlab.com/kafa1942/take10dashboard/common"
	"gitlab.com/kafa1942/take10dashboard/docs"
	"gitlab.com/kafa1942/take10dashboard/users"
)

// migration function, used to migrate modules models if needed
func Migrate(db *gorm.DB) {
	users.AutoMigrate()
}

func main() {
	
	// swagger documentation information, to change and see 
	// other configuration see doc.go 
	docs.SwaggerInfo.Title = common.Config.App.Name
	docs.SwaggerInfo.Description = common.Config.App.Description
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = common.Config.App.Host
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// init database and Migrate if needed
	db := common.GetDB()
	Migrate(db)
	defer db.Close()

	// init gin router
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		users.UsersRegister(v1.Group("/users"))
		v1.Use(users.AuthMiddleware(true))
		users.UserRegister(v1.Group("/user"))
	}
	
	// health check
	healthCheck := r.Group("/api/ping")

	healthCheck.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// swagger docs files , TODO only for staging and development environments, add condition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	_ = r.Run() // listen and serve on 0.0.0.0:8080
}
