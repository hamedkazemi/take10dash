package promotions

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/kafa1942/take10dashboard/common"
	_ "gitlab.com/kafa1942/take10dashboard/users"
	"net/http"
)

func GetAll(c *gin.Context) {
	var input common.GetAllRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Error("Promotions", err)
		c.JSON(http.StatusBadRequest, common.NewError("Promotions", errors.New("Bad Request Input")))
		return
	}
	var searchQuery string
	var orderBy string
	var orderType string

	limit := input.Limit
	offset := input.Offset

	if input.Query != "" {
		searchQuery = input.Query
	}
	if input.OrderBy != "" {
		orderBy = input.OrderBy
	} else {
		orderBy = "id"
	}
	if input.OrderType != "" {
		orderType = input.OrderType
	} else {
		orderType = "DESC"
	}

	filters := input.Filters

	if searchQuery != "" {
		filters = append(filters, common.Filter{Field: "allKeys", Operator: "LIKE", Value: searchQuery})
	}

	model, modelCount, err := FindMany(filters, orderBy, orderType, limit, offset)
	if err != nil {
		logrus.Error("Promotions", err)
		c.JSON(http.StatusNotFound, common.NewError("Promotions", errors.New("Something went wrong!")))
		return
	}

	var metaData common.MetaData
	metaData.Pagination.Count = modelCount
	metaData.Pagination.Limit = limit
	metaData.Pagination.Offset = offset
	metaData.Filters = filters
	metaData.Order.OrderBy = orderBy
	metaData.Order.OrderType = orderType

	c.JSON(http.StatusOK, gin.H{"data": model, "metaData": metaData})
}

func Get(c *gin.Context) {
	var r GetRequest
	if err := c.ShouldBindUri(&r); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("message", errors.New(err.Error())))
		return
	}
	var model Promotion

	nf := common.GetDB().First(&model, r.ID).RecordNotFound()

	if !nf {
		c.JSON(http.StatusOK, gin.H{"data": model})
	} else {
		c.JSON(http.StatusNotFound, common.NewError("Promotions", errors.New("Not Found.")))
	}

}

func Update(c *gin.Context) {
	var input PromotionUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Error("Promotions", err)
		c.JSON(http.StatusBadRequest, common.NewError("Promotions", errors.New("Bad Request Input")))
		return
	}
	var model Promotion

	nf := common.GetDB().First(&model, input.ID).RecordNotFound()

	if !nf {
		model.Status = input.Status
		model.Code = input.Code
		model.EndDate = input.EndDate
		model.StartDate = input.StartDate
		model.Points = input.Points
		model.PrizeURL = input.PrizeURL
		model.PromoCategory = input.PromoCategory

		err := common.GetDB().Save(&model).Error
		if err != nil {
			logrus.Error("Promotions", err)
			c.JSON(http.StatusNotFound, common.NewError("Promotions", errors.New("Can't Update Right Now.")))
		} else {
			c.JSON(http.StatusOK, gin.H{"data": model})
		}
	} else {
		c.JSON(http.StatusNotFound, common.NewError("Promotions", errors.New("Not Found.")))
	}

}

func Delete(c *gin.Context) {
	var r GetRequest
	if err := c.ShouldBindUri(&r); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("message", errors.New(err.Error())))
		return
	}
	var model Promotion

	nf := common.GetDB().First(&model, r.ID).RecordNotFound()

	if !nf {
		err := common.GetDB().Delete(&model).Error
		if err != nil {
			logrus.Error("Promotions", err)
			c.JSON(http.StatusNotFound, common.NewError("Promotions", errors.New("Can't Delete Right Now.")))
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Promotion Deleted Successfully"})
		}
	} else {
		c.JSON(http.StatusNotFound, common.NewError("Promotions", errors.New("Not Found.")))
	}
}

func Create(c *gin.Context) {
	var input PromotionCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Error("Promotions", err)
		c.JSON(http.StatusBadRequest, common.NewError("Promotions", errors.New("Bad Request Input")))
		return
	}
	var model Promotion
	model.Status = input.Status
	model.Code = input.Code
	model.EndDate = input.EndDate
	model.StartDate = input.StartDate
	model.Points = input.Points
	model.PrizeURL = input.PrizeURL
	model.PromoCategory = input.PromoCategory

	err := common.GetDB().Create(&model).Error
	if err != nil {
		logrus.Error("Promotions", err)
		c.JSON(http.StatusConflict, common.NewError("Promotions", errors.New("Duplicate Entry, Can't Create Right Now.")))
	} else {
		c.JSON(http.StatusOK, gin.H{"data": model, "message": "Promotion Created Successfully"})
	}
}
