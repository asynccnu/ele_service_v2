package service

import (
	"github.com/asynccnu/ele_service_v2/util"
)

type RemainPowerPayload struct {
	Remain   string `json:"remain" xml:"remainPower"` // 剩余电量，单位：度
	ReadTime string `json:"read_time" xml:"readTime"` // 最近一次抄表时间
}

// 请求获取剩余电量（度）
func MakeRemainPowerRequest(meterId string) (*RemainPowerPayload, error) {
	url := "https://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/getReserveHKAM?AmMeter_ID=" + meterId

	var payload RemainPowerPayload
	if err := MakeRequest(url, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

// 请求用电详情解析结构
type DayValueInfo struct {
	List []*DayElectricItem `xml:"dayValueInfoList>DayValueInfo"`
}

// 每日用电信息
type DayElectricItem struct {
	Usage  string `xml:"dayValue"`    // 用电量
	Charge string `xml:"dayUseMeony"` // 电费
}

// 获取用电详情，开始时间-结束时间
func MakeChargeInfoRequest(meterId string, startDate, endDate string) ([]*DayElectricItem, error) {
	url := "https://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/getMeterDayValue?AmMeter_ID=" + meterId +
		"&startDate=" + startDate + "&endDate=" + endDate

	var payload DayValueInfo
	if err := MakeRequest(url, &payload); err != nil {
		return nil, err
	}

	return payload.List, nil
}

// 发起请求，解析数据
func MakeRequest(url string, data interface{}) error {
	// 发送 HTTP GET 请求
	body, err := util.SendHTTPGetRequest(url)
	if err != nil {
		return err
	}

	// 解析 XML body data
	if err := util.UnmarshalXMLBody(body, data); err != nil {
		return err
	}

	return nil
}
