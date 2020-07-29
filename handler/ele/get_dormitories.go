package ele

import (
	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/service"
	"github.com/gin-gonic/gin"
	"strings"
)

// 获取房间信息
func GetDormitories(c *gin.Context) {
	architecture := c.Query("architecture_id")
	floor := c.Query("floor")
	if architecture == "" || floor == "" {
		handler.SendError(c, nil, nil, "missing parameters")
		return
	}

	dormitories, err := service.GetRoomId(architecture, floor)
	if err != nil {
		log.Error("GetRoomId function error")
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// 传给前端的宿舍名去掉空调照明后缀,还要去重
	var trimDorms []string
	store := make(map[string]bool)

	for i := 0; i < len(dormitories.Dorms); i++ {
		s := strings.Trim(dormitories.Dorms[i].DormName, "空调照明")
		if _, exist := store[s]; !exist {
			trimDorms = append(trimDorms, s)
		}
		store[s] = true
	}

	handler.SendResponse(c, nil, trimDorms)
}
