package ele

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/model"
	"github.com/asynccnu/ele_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

type GetDormitoryResponse struct {
	Count int      `json:"count"`
	List  []string `json:"list"`
}

// 获取宿舍
func GetDormitories(c *gin.Context) {
	building := c.DefaultQuery("building", "")
	if building == "" {
		handler.SendBadRequest(c, errno.ErrQuery, nil, "building is required.")
		return
	}

	list, err := model.GetRoomsByBuildingName(building)
	if err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// 东7需另外考虑，有 101/101A/101B 三个宿舍号
	// 选择宿舍号需要排除 101
	if building == "东7" {
		for i := len(list) - 1; i >= 0; i-- {
			if !strings.HasSuffix(list[i], "A") && !strings.HasSuffix(list[i], "B") {
				fmt.Println("1")
				list = append(list[:i], list[i+1:]...)
			}
		}
	}

	sort.Strings(list)

	handler.SendResponse(c, nil, &GetDormitoryResponse{
		Count: len(list),
		List:  list,
	})
}
