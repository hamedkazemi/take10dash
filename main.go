package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"gitlab.com/kafa1942/take10dashboard/common"
	"gitlab.com/kafa1942/take10dashboard/docs"
	"gitlab.com/kafa1942/take10dashboard/users"
	"net"
	"net/http"
	"net/url"
	"os"
	"take10dashboard/gameusers"
	"take10dashboard/promotions"
	"take10dashboard/questions"
	"time"
)

// migration function, used to migrate modules models if needed
func Migrate(db *gorm.DB) {
	users.AutoMigrate()
}

func main() {
	logrus.Info("API system loading ...")

	// swagger documentation information, to change and see
	// other configuration see doc.go
	docs.SwaggerInfo.Title = common.Config.App.Name
	docs.SwaggerInfo.Description = common.Config.App.Description
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = common.Config.App.Host
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// set proxy for outgoing requests ( if needed )
	setProxy()

	// use kafka client
	//common.KafkaCon()

	// init database and Migrate if needed
	db := common.GetDB()
	logrus.Info("Migrating db if needed...")
	Migrate(db)
	defer db.Close()

	// init gin router
	r := gin.Default()
	r.Use(common.CORS())
	v1 := r.Group("/api/v1")
	{
		questions.ConfigGinRouter(v1)
		users.ConfigGinRouter(v1)
		gameusers.ConfigGinRouter(v1)
		promotions.ConfigGinRouter(v1)
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

	logrus.Info("APP will be served at: " + common.Config.App.Host)
	_ = r.Run(common.Config.App.Host)
}

func setProxy() {
	getProxy := ""
	if os.Getenv("HTTP_PROXY") != "" {
		getProxy = os.Getenv("HTTP_PROXY")
	} else {
		getProxy = common.Config.App.Proxy
	}

	proxyUrl, err := url.Parse(getProxy)
	if err != nil {
		logrus.Error("the URL of Proxy is wrong.")
		logrus.Panic("the URL of Proxy is wrong, please check config.toml file.")
	}

	http.DefaultTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy:                 http.ProxyURL(proxyUrl),
	}
}
