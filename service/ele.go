package service

import (
	"errors"
	"strings"
	"time"

	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/model"
)

// 获取剩余电量
func GetRemainPower(meterId string) (*RemainPowerPayload, error) {
	return MakeRemainPowerRequest(meterId)
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

// 获取电费，剩余电量和昨日用量
func GetElectricChargeByMeterID(meterId string) (*model.ElectricityInfo, error) {
	remainPower, err := GetRemainPower(meterId)
	if err != nil {
		log.Error("GetRemainPower function error")
		return nil, err
	}

	yesterdayCharge, err := GetYesterdayCharge(meterId)
	if err != nil {
		log.Error("GetYesterdayCharge function error")
		return nil, err
	}

	var data = &model.ElectricityInfo{
		RemainPower: remainPower.Remain,
		ReadTime:    remainPower.ReadTime,
		Consumption: model.EleConsumptionItem{
			Usage:  yesterdayCharge.Usage,
			Charge: yesterdayCharge.Charge,
		},
	}

	return data, nil
}

// 根据楼栋和宿舍号获取电费
func GetElectricity(building, room string) ([]*model.ElectricityInfo, error) {
	var list = make([]*model.ElectricityInfo, 0)

	meters, err := model.GetMetersByBuildingAndRoom(building, room)
	if err != nil {
		log.Error("GetMetersByBuildingAndRoom function error: " + err.Error())
		return nil, err
	}

	for _, meter := range meters {
		// 获取电费
		record, err := GetElectricChargeByMeterID(meter.MeterID)
		if err != nil {
			log.Error("GetElectricChargeByMeterID function error: " + err.Error())
			return nil, err
		}
		record.Kind = meter.Kind

		list = append(list, record)
	}

	// 东7宿舍号为 101A 或 101B
	// 获取照明和客厅空调电费需要去除后缀的 A/B 来请求
	if building == "东7" && (strings.Contains(room, "A") || strings.Contains(room, "B")) {
		room := strings.TrimSuffix(strings.TrimSuffix(room, "B"), "A")
		meters, err := model.GetMetersByBuildingAndRoom(building, room)
		if err != nil {
			log.Error("GetMetersByBuildingAndRoom function error: " + err.Error())
			return nil, err
		}

		for _, meter := range meters {
			record, err := GetElectricChargeByMeterID(meter.MeterID)
			if err != nil {
				log.Error("GetElectricChargeByMeterID function error: " + err.Error())
				return nil, err
			}

			// 101 室的照明是 101A、101B 和 客厅共用的，可直接并作是 101A/101B 的照明
			// 空调是客厅的，101A/101B 都有各自的室内空调，因此需要另外表示
			if meter.Kind == "air" {
				record.Kind = "客厅-空调"
			} else {
				record.Kind = meter.Kind
			}

			list = append(list, record)
		}
	}

	return list, nil
}
