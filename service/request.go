package service

import (
	"github.com/asynccnu/ele_service_v2/util"
)

// 请求获取楼栋信息
func MakeArchitecureRequest(areaId string) (*ArchitectureList, error) {
	var result ArchitectureList
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/getArchitectureInfo?Area_ID=" + areaId

	if err := MakeRequest(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 请求获取寝室信息
func MakeDormitoryRequest(architectureId string, floorId string) (*DormList, error) {
	var result DormList
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getRoomInfo?Architecture_ID=" + architectureId + "&Floor=" + floorId

	if err := MakeRequest(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type MeterInfoPayload struct {
	MeterId string `xml:"meterList>MeterInfo>meterId"` // 电表号 id
}

// 请求获取电表信息
func MakeMeterInfoRequest(roomId string) (string, error) {
	var data MeterInfoPayload
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/getMeterInfo?Room_ID=" + roomId

	if err := MakeRequest(url, &data); err != nil {
		return "", err
	}

	return data.MeterId, nil
}

type RemainPowerPayload struct {
	Remain   string `json:"remain" xml:"remainPower"` // 剩余电量，单位：度
	ReadTime string `json:"read_time" xml:"readTime"` // 最近一次抄表时间
}

// 请求获取剩余电量（度）
func MakeRemainPowerRequest(meterId string) (*RemainPowerPayload, error) {
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/getReserveHKAM?AmMeter_ID=" + meterId
	var payload RemainPowerPayload
	if err := MakeRequest(url, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

type DayValueInfo struct {
	List []*DayElectricItem `xml:"dayValueInfoList>DayValueInfo"`
}

type DayElectricItem struct {
	Use string `xml:"dayValue"`    // 昨日用电量
	Fee string `xml:"dayUseMeony"` // 昨日电费
}

// 获取用电详情，开始时间-结束时间
func MakeChargeInfoRequest(meterId string, startDate, endDate string) ([]*DayElectricItem, error) {
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/getMeterDayValue?AmMeter_ID=" + meterId +
		"&startDate=" + startDate + "&endDate=" + endDate

	var payload DayValueInfo
	if err := MakeRequest(url, &payload); err != nil {
		return nil, err
	}

	return payload.List, nil
}

// 发起请求，解析数据
func MakeRequest(url string, data interface{}) error {
	body, err := util.SendHTTPGetRequest(url)
	if err != nil {
		return err
	}

	if err = util.UnmarshalXMLBody(body, &data); err != nil {
		return err
	}
	return nil
}
