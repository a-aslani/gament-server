package controllerV1

import (
	"Server/app/model"
	"Server/app/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//Create Role
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
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "create new rule failed"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.CreateDocument{
			Message: "Rule successfully created",
			Key:     ruleID.Key(),
		},
		State: true,
	})
}
