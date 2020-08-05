package service

import (
	"errors"
	"time"

	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/model"
)

var AreaMap = map[string]string{"1": "0001", "2": "0002", "3": "0003", "4": "0004", "5": "0006"}

// 获取每个区域楼层信息
// 西区学生宿舍,东区学生宿舍,元宝山学生宿舍,南湖学生宿舍,国际园区 分别是1 2 3 4 5.
func GetArchitectures(area string) (*ArchitectureList, error) {
	areaId, ok := AreaMap[area]
	if !ok {
		return nil, errors.New("area is error")
	}

	return MakeArchitecureRequest(areaId)
}

// 获取寝室信息
func GetRoomId(architectureId string, floorId string) (*DormList, error) {
	return MakeDormitoryRequest(architectureId, floorId)
}

//  获取电表数据
func GetMeterInfo(area, architecture, floor, dormName string) (*model.MeterInfo, error) {
	var (
		architectureId string
		dormInfo       DormInfo
	)

	// 获取楼栋ID
	architecures, err := MakeArchitecureRequest(area)
	if err != nil {
		return nil, err
	}

	for _, item := range architecures.List {
		if item.ArchitectureName == architecture {
			architectureId = item.ArchitectureID
			break
		}
	}

	// 获取宿舍ID
	dormList, err := MakeDormitoryRequest(architectureId, floor)
	if err != nil {
		return nil, err
	}

	for _, item := range dormList.Dorms {
		if item.DormName == dormName {
			dormInfo.DormId = item.DormId
			// dormInfo.DormName = dormName
			break
		}
	}

	// 获取电表ID
	meterId, err := MakeMeterInfoRequest(dormInfo.DormId)
	if err != nil {
		return nil, err
	}

	var meterInfo = &model.MeterInfo{
		DormName: dormName,
		MeterId:  meterId,
	}

	// 如果能获取到电表信息,就把电表信息添加到数据库中
	if err = model.AddMeterInfo(meterInfo); err != nil {
		log.Error(" AddMeterInfo function error " + err.Error())
	}

	return meterInfo, nil
}

// 获取剩余电量
func GetRemainPower(meterId string) (*RemainPowerPayload, error) {
	remainPower, err := MakeRemainPowerRequest(meterId)
	if err != nil {
		return nil, err
	}
	return remainPower, nil
}

// 昨日用电量
func GetYesterdayCharge(meterId string) (*DayElectricItem, error) {
	date := time.Now().AddDate(0, 0, -1).Format("2006/01/02")
	chargeList, err := MakeChargeInfoRequest(meterId, date, date)
	if err != nil {
		return nil, err
	}

	if len(chargeList) < 1 {
		return nil, errors.New("")
	}

	return chargeList[0], nil
}

// 有两个请求 第一个请求获取电费和时间，第二个请求获取昨日用量
func GetElectricChargeByMeterId(meterId string) (*model.ElectricCharge, error) {
	remainPower, err := GetRemainPower(meterId)
	if err != nil {
		log.Error("GetRemainPower function error")
	}

	yesterdayCharge, err := GetYesterdayCharge(meterId)
	if err != nil {
		log.Error("GetYesterdayCharge function error")
	}

	var charge = model.ElectricCharge{
		Id: meterId,
		// Type:        meterType,
		RemainPower: remainPower.Remain,
		ReadTime:    remainPower.ReadTime,
		Yesterday: model.DayElectricInfo{
			DayElecUse: yesterdayCharge.Use,
			DayFee:     yesterdayCharge.Fee,
		},
	}

	// if err != nil {
	// 	cache, err := model.GetElectricity(meterId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return cache, nil
	// }

	// if err := model.AddElectricity(charge); err != nil {
	// 	log.Error("AddElectricity function error" + err.Error())
	// }

	return &charge, nil
}

type MeterRequestPayload struct {
	Area         string `json:"area" binding:"true"`
	Architecture string `json:"architecture" binding:"true"`
	Floor        string `json:"floor" binding:"true"`
	Dormitory    string `json:"dormitory" binding:"true"` // 东16-233
}

func GetElectricity(meterType string, payload *MeterRequestPayload) (*model.ElectricCharge, error) {
	dormName := payload.Dormitory
	if meterType == "light" {
		dormName += "照明"
	} else if meterType == "air" {
		dormName += "空调"
	} else {
		return nil, nil
	}

	meterInfo, err := GetMeterInfo(payload.Area, payload.Architecture, payload.Floor, dormName)
	if err != nil {
		return nil, nil
	}

	data, err := GetElectricChargeByMeterId(meterInfo.MeterId)
	if err != nil {
		return nil, nil
	}
	data.Type = meterType

	return data, nil
}
