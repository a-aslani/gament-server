package controllerApiV1

import (
	"Server/app/model"
	"Server/app/response"
	"Server/app/utility"
	"github.com/arangodb/go-driver"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//Create Rule
func CreateRule(c *gin.Context) {

	var rule model.Rule
	if err := c.ShouldBind(&rule); err != nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: err.Error()}})
		return
	}

	t := time.Now().Unix()
	rule.CreatedAt = t
	rule.UpdatedAt = t

	ruleID, err := model.Create(&rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "مشکل در ایجاد قوانین جدید"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.CreateDocument{
			Message: "قوانین جدید با موفقیت ایجاد شد",
			Key:     ruleID.Key(),
		},
		State: true,
	})
}

//Find rule by key
func FindRule(c *gin.Context) {

	ruleKey := c.Param("key")

	rule, err := model.Find(ruleKey, &model.Rule{})

	if driver.IsNotFound(err) {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "قانونی پیدا نشد"}})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "خطا در پیدا کردن قانون مورد نظر"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.FindDocument{
			Document: utility.RefactorResponseDoc(rule, false),
		},
		State: true,
	})
}
