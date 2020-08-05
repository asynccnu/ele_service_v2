package ele

import (
	"strings"

	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/pkg/errno"
	"github.com/asynccnu/ele_service_v2/service"

	"github.com/gin-gonic/gin"
)

type GetDormitoryResponse struct {
	Count int      `json:"count"`
	List  []string `json:"list"`
}

// 获取房间信息
func GetDormitories(c *gin.Context) {
	architecture := c.Query("architecture")
	floor := c.Query("floor")
	if architecture == "" || floor == "" {
		handler.SendBadRequest(c, errno.ErrQuery, nil, "architecture and floor are required.")
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
			store[s] = true
		}
	}

	handler.SendResponse(c, nil, &GetDormitoryResponse{
		Count: len(trimDorms),
		List:  trimDorms,
	})
}
