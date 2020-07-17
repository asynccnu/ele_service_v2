package ele

import (
	"errors"
	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/model"
	"github.com/asynccnu/ele_service_v2/service"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	area := c.Query("area")
	architecture := c.Query("architecture")
	floor := c.Query("floor")
	dormitory := c.Query("dormitory")

	var meterInfo model.MeterInfo
	var err error
	empty := model.MeterInfo{}

	//获取电表信息
	if meterInfo, err = model.GetMongoMeterInfo(dormitory); err != nil || meterInfo == empty {
		if meterInfo, err = service.GetMeterInfoAPI(area, architecture, floor, dormitory); err != nil {
			log.Error("GetMongoMeterInfo function error")
			handler.SendError(c, err, nil, err.Error())
			return
		}
		if meterInfo == empty {
			log.Error("The request parameter is out of specification")
			err := errors.New("the request parameter is out of specification")
			handler.SendError(c, err, nil, "")
			return
		}
		//把电表信息添加到数据库中
		if err = model.AddMeterInfo(meterInfo); err != nil {
			log.Error(" AddMeterInfo function error")
		}
	}

	//获取电费信息
	charge, err := service.GetElectricCharge(meterInfo.MeterId)
	if err != nil {
		log.Error("GetElectricityCharge function error")
		handler.SendError(c, err, nil, err.Error())
		return
	}
	handler.SendResponse(c, nil, charge)

}
