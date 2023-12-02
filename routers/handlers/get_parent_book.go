package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/types"
	"net/http"
)

func GetParentBookInfo(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.PureJSON(http.StatusBadRequest, "not set id param")
		return
	}
	info, err := types.GetBookGroupInfoByChildBookID(id)
	if err != nil {
		logger.Info(err)
		c.PureJSON(http.StatusBadRequest, "ParentBookInfo not found")
		return
	}
	c.PureJSON(http.StatusOK, info)
}
