package ele

import (
	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/model"
	"github.com/asynccnu/ele_service_v2/service"
	"github.com/gin-gonic/gin"
	"strings"
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

// 获取电表信息
func Get(c *gin.Context) {
	// 获取电表信息
	area := c.Query("area")
	architecture := c.Query("architecture")
	floor := c.Query("floor")
	dormitory := c.Query("dormitory")
	meterType := c.Query("type")

	var meterInfoList []model.MeterInfo
	var dormitories []string
	var chargeList []model.ElectricCharge
	var meterInfo model.MeterInfo
	var err error
	empty := model.MeterInfo{}

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
		dormitories = append(dormitories, dormitory+"空调")
		dormitories = append(dormitories, dormitory+"照明")
	}
	dormitories = append(dormitories, dormitory)

	// 获取电表信息
	for i := 0; i < len(dormitories); i++ {
		if meterInfo, err = model.GetMongoMeterInfo(dormitories[i]); err != nil || meterInfo == empty {
			// 参数错的话不会报错 比如原来是不用加空调照明后缀的,我全部加了也不会报错
			if meterInfo, err = service.GetMeterInfoAPI(area, architecture, floor, dormitories[i]); err != nil {
				log.Error("GetMongoMeterInfo function error" + err.Error())
				handler.SendError(c, err, nil, err.Error())
				return
			}

			// 如果能获取到电表信息,就把电表信息添加到数据库中
			if meterInfo != empty {
				if err = model.AddMeterInfo(meterInfo); err != nil {
					log.Error(" AddMeterInfo function error" + err.Error())
				}
			}
		}
		if meterInfo != empty {
			meterInfoList = append(meterInfoList, meterInfo)
		}
	}

	if len(meterInfoList) == 0 {
		log.Error("Can't get right meterInfo")
		handler.SendError(c, nil, nil, "may be wrong request parameters")
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
		chargeList = append(chargeList, charge)
	}

	handler.SendResponse(c, nil, chargeList)
}
