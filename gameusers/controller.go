package gameusers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/kafa1942/take10dashboard/common"
	_ "gitlab.com/kafa1942/take10dashboard/users"
	"net/http"
)

//router.GET("/", GetAll)
//router.GET("/:id", Get)
//router.PATCH("/", Update)
//router.DELETE("/", Delete)

func GetAll(c *gin.Context) {
	var input common.GetAllRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Error("Users", err)
		c.JSON(http.StatusBadRequest, common.NewError("Users", errors.New("Bad Request Input")))
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
		orderBy = "created_at"
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
		logrus.Error("Users", err)
		c.JSON(http.StatusNotFound, common.NewError("Users", errors.New("Something went wrong!")))
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
	var model User

	nf := common.GetDB().First(&model, r.ID).RecordNotFound()

	if !nf {
		c.JSON(http.StatusOK, gin.H{"data": model})
	} else {
		c.JSON(http.StatusNotFound, common.NewError("Users", errors.New("Not Found.")))
	}

}

func Update(c *gin.Context) {
	var input UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Error("Users", err)
		c.JSON(http.StatusBadRequest, common.NewError("Users", errors.New("Bad Request Input")))
		return
	}
	var model User

	nf := common.GetDB().First(&model, input.UserID).RecordNotFound()

	if !nf {
		model.Email = input.Email
		model.FbFirstName = input.FbFirstName
		model.FbLastName = input.FbLastName
		model.PhoneNumber = input.PhoneNumber
		model.Name = input.Name
		model.UserStatus = input.UserStatus
		err := common.GetDB().Save(&model).Error
		if err != nil {
			logrus.Error("Users", err)
			c.JSON(http.StatusNotFound, common.NewError("Users", errors.New("Can't Update Right Now.")))
		} else {
			c.JSON(http.StatusOK, gin.H{"data": model})
		}
	} else {
		c.JSON(http.StatusNotFound, common.NewError("Users", errors.New("Not Found.")))
	}

}

func Delete(c *gin.Context) {
	var r GetRequest
	if err := c.ShouldBindUri(&r); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("message", errors.New(err.Error())))
		return
	}
	var model User

	nf := common.GetDB().First(&model, r.ID).RecordNotFound()

	if !nf {
		err := common.GetDB().Delete(&model).Error
		if err != nil {
			logrus.Error("Users", err)
			c.JSON(http.StatusNotFound, common.NewError("Users", errors.New("Can't Delete Right Now.")))
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "User Deleted Successfully"})
		}
	} else {
		c.JSON(http.StatusNotFound, common.NewError("Users", errors.New("Not Found.")))
	}
}

func GetProfile(c *gin.Context) {
	var r GetRequest
	if err := c.ShouldBindUri(&r); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("message", errors.New(err.Error())))
		return
	}
	var UserModel User
	var UserProfile UserProfile
	var phoneVerification PhoneVerification
	var userActivity ActivityLog
	var userActivities []ActivityLog
	var Winnings []RandomWinner

	var resp ProfileResponse

	nf := common.GetDB().First(&UserModel, r.ID).RecordNotFound()

	if !nf {
		resp.User = UserModel
	} else {
		c.JSON(http.StatusNotFound, common.NewError("Users", errors.New("Not Found.")))
		return
	}

	nf = common.GetDB().First(&UserProfile, r.ID).RecordNotFound()

	if !nf {
		resp.UserProfile = UserProfile
	} else {
		c.JSON(http.StatusNotFound, common.NewError("Users", errors.New("Not Found.")))
		return
	}

	var invitePromoCount int
	common.GetDB().Model(InvitePromo{}).Where("sender_user_id = ?", r.ID).Count(&invitePromoCount)
	resp.InvitePromoCount = invitePromoCount

	var totalCompletePayment string
	_ = common.GetDB().Model(CompletePayment{}).Select("sum(amount) as total").Where("user_id = ?", r.ID).Row().Scan(&totalCompletePayment)
	resp.CompletePayment = totalCompletePayment

	common.GetDB().Where("user_id = ?", r.ID).Order("created_at DESC").First(&phoneVerification)
	resp.LastPhoneVerificationDate = phoneVerification.CreatedAt
	resp.LastPhoneVerificationStatus = phoneVerification.Status

	common.GetDB().Where("activity_message = 'user_logged_in_success' AND user_id = ?", r.ID).Order("created_at DESC").First(&userActivity)
	resp.LastLoggedInDate = userActivity.CreatedAt

	common.GetDB().Where("user_id = ?", r.ID).Order("created_at DESC").Limit("20").Find(&userActivities)
	resp.LastActivities = userActivities

	common.GetDB().Where("user_id = ?", r.ID).Order("date DESC").Find(&Winnings)
	resp.Winnings = Winnings

	c.JSON(http.StatusOK, gin.H{"data": resp})

}
