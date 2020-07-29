package ele

import (
	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/service"
	"github.com/gin-gonic/gin"
)

// 获取楼栋信息
func GetArchitectures(c *gin.Context) {
	area := c.Query("area")
	if area == "" {
		handler.SendError(c, nil, nil, "missing parameter area")
		return
	}

	architectures, err := service.GetArchitectures(area)
	if err != nil {
		log.Error("GetArchitectures function error")
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, &architectures)

}
