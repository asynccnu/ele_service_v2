package ele

import (
	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/model"
	"github.com/asynccnu/ele_service_v2/pkg/errno"
	"github.com/asynccnu/ele_service_v2/service"

	"github.com/gin-gonic/gin"
)

type GetElectricityResponse struct {
	Building string `json:"building"`
	Room     string `json:"room"`

	HasLight bool                   `json:"has_light"`
	Light    *model.ElectricityInfo `json:"light"`
	HasAir   bool                   `json:"has_air"`
	Air      *model.ElectricityInfo `json:"air"`

	// 一些宿舍楼可能一间宿舍会有其它的数据，
	// 比如东7的宿舍还有客厅，要负责客厅的空调电费
	HasMore  bool                     `json:"has_more"`
	MoreData []*model.ElectricityInfo `json:"more_data"`
}

// 获取电表信息
func Get(c *gin.Context) {
	building := c.DefaultQuery("building", "")
	room := c.DefaultQuery("room", "")
	if building == "" || room == "" {
		handler.SendBadRequest(c, errno.ErrQuery, nil, "building and room are required.")
		return
	}

	records, err := service.GetElectricity(building, room)
	if err != nil {
		handler.SendError(c, errno.ErrGetEle, nil, err.Error())
		return
	}

	var responseData = &GetElectricityResponse{
		Building: building,
		Room:     room,
	}

	for _, record := range records {
		switch record.Kind {
		case "air":
			responseData.HasAir = true
			responseData.Air = record
		case "light":
			responseData.HasLight = true
			responseData.Light = record
		default:
			responseData.HasMore = true
			responseData.MoreData = append(responseData.MoreData, record)
		}
	}

	handler.SendResponse(c, nil, responseData)
}
