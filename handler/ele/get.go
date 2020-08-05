package ele

import (
	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/pkg/errno"
	"github.com/asynccnu/ele_service_v2/service"

	"github.com/gin-gonic/gin"
)

// 获取电表信息
func Get(c *gin.Context) {
	// area := c.Query("area")
	// architecture := c.Query("architecture")
	// floor := c.Query("floor")
	// dormitory := c.Query("dormitory")
	// if area == "" || architecture == "" || floor == "" || dormitory == "" {
	// 	handler.SendBadRequest(c, errno.ErrQuery, nil, "area, architecture, floor and dormitory are all required.")
	// 	return
	// }

	meterType := c.DefaultQuery("type", "light")

	payload := &service.MeterRequestPayload{
		// Area:      area,
		// Building:  architecture,
		// Floor:     floor,
		// Dormitory: dormitory,
	}

	if err := c.BindQuery(payload); err != nil {
		handler.SendBadRequest(c, errno.ErrQuery, nil, "query")
		return
	}

	data, err := service.GetElectricity(meterType, payload)
	if err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, "")
		return
	}

	handler.SendResponse(c, nil, data)
}
