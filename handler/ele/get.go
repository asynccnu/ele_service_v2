package ele

import (
	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/model"
	"github.com/asynccnu/ele_service_v2/service"
	"github.com/gin-gonic/gin"
)

// 获取电表信息
func Get(c *gin.Context) {
	// 获取电表信息 缺少参数和参数错误的情况在下面判断
	area := c.Query("area")
	architecture := c.Query("architecture")
	floor := c.Query("floor")
	dormitory := c.Query("dormitory")
	meterType := c.DefaultQuery("type", "all")

	var (
		meterInfoList []*model.MeterInfo
		dormitories   []string
		chargeList    []model.ElectricCharge
		meterInfo     *model.MeterInfo
		err           error
	)
	switch meterType {
	case "air":
		dormitory += "空调"
	case "light":
		dormitory += "照明"
	default:
		meterType = "all"
	}

	// 有的国交南湖宿舍是不分空调和照明的,就都放进去请求
	if meterType == "all" {
		dormitories = append(dormitories, dormitory+"空调", dormitory+"照明")
	}
	dormitories = append(dormitories, dormitory)

	// 获取电表信息
	for _, v := range dormitories {
		if meterInfo, err = model.GetMongoMeterInfo(v); err != nil || *meterInfo == (model.MeterInfo{}) {
			// 参数错的话不会报错 比如原来是不用加空调照明后缀的,我全部加了也不会报错
			if meterInfo, err = service.GetMeterInfoAPI(area, architecture, floor, v); err != nil {
				log.Error("GetMeterInfoAPI function error" + err.Error())
				handler.SendError(c, err, nil, err.Error())
				return
			}

			// 如果能获取到电表信息,就把电表信息添加到数据库中
			if *meterInfo != (model.MeterInfo{}) {
				if err = model.AddMeterInfo(meterInfo); err != nil {
					log.Error(" AddMeterInfo function error " + err.Error())
				}
			}
		}
		if *meterInfo != (model.MeterInfo{}) {
			meterInfoList = append(meterInfoList, meterInfo)
		}
	}

	// 参数错误和缺少参数归在一类
	if len(meterInfoList) == 0 {
		log.Error("Can't get right meterInfo")
		handler.SendError(c, nil, nil, "may be wrong request parameters or missing parameters")
		return
	}

	// 获取电费信息
	for i := 0; i < len(meterInfoList); i++ {
		if len(meterInfoList) == 2 {
			if i == 0 {
				meterType = "air"
			} else {
				meterType = "light"
			}
		}

		charge, err := service.GetElectricCharge(meterInfoList[i].MeterId, meterType)
		if err != nil {
			log.Error("GetElectricityCharge function error")
			handler.SendError(c, err, nil, err.Error())
			return
		}

		chargeList = append(chargeList, *charge)
	}

	handler.SendResponse(c, nil, chargeList)
}
