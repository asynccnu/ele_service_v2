package ele

import (
	"errors"
	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/model"
	"github.com/asynccnu/ele_service_v2/service"
	"github.com/gin-gonic/gin"
)

// 获取楼栋信息
func GetArchitectures(c *gin.Context) {
	area := c.Query("area")

	architectures, err := service.GetArchitectures(area)
	if err != nil {
		log.Error("GetArchitectures function error")
		handler.SendError(c, err, nil, err.Error())
		return
	}
	handler.SendResponse(c, nil, architectures)

}

// 获取房间信息
func GetDormitories(c *gin.Context) {
	architecture := c.Query("architecture_id")
	floor := c.Query("floor")

	dormitories, err := service.GetRoomId(architecture, floor)
	if err != nil {
		log.Error("GetRoomId function error")
		handler.SendError(c, err, nil, err.Error())
		return
	}
	handler.SendResponse(c, nil, dormitories)
}

// 获取电表信息
func Get(c *gin.Context) {
	// 获取电表信息
	dormName := c.Query("dorm_name")
	dormId := c.Query("dorm_id")
	var meterInfo model.MeterInfo
	var err error
	empty := model.MeterInfo{}

	//获取电表信息
	if meterInfo, err = model.GetMongoMeterInfo(dormName); err != nil || meterInfo == empty {
		if meterInfo, err = service.GetMeterInfo(dormId); err != nil {
			log.Error("GetMongoMeterInfo and GetMeterInfo function error")
			handler.SendError(c, err, nil, err.Error())
			return
		}
		if meterInfo == empty {
			log.Error("The request parameter is out of specification")
			err := errors.New("the request parameter is out of specification")
			handler.SendError(c, err, nil, "")
			return
		}
		meterInfo.DormName = dormName
		//把电表信息添加到数据库中
		if err = model.AddMeterInfo(meterInfo); err != nil {
			log.Error(" AddMeterInfo function error")
		}
	}

	// 获取电费信息
	charge, err := service.GetElectricCharge(meterInfo.MeterId)
	if err != nil {
		log.Error("GetElectricityCharge function error")
		handler.SendError(c, err, nil, err.Error())
		return
	}
	handler.SendResponse(c, nil, charge)

}
