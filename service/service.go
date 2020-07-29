package service

import (
	"encoding/xml"
	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/model"
	"io/ioutil"
	"net/http"
	"time"
)

// 宿舍信息
type DormList struct {
	Dorms []DormInfo `xml:"roomInfoList>RoomInfo"`
}

// 宿舍信息
type DormInfo struct {
	DormName string `xml:"RoomName"`
	DormId   string `xml:"RoomNo" `
}

// 楼栋信息
type ArchiList struct {
	Architectures []Architecture `xml:"architectureInfoList>architectureInfo" json:"architectures"`
}

// 楼栋信息
type Architecture struct {
	ArchitectureID   string `xml:"ArchitectureID" json:"architecture_id"`
	ArchitectureName string `xml:"ArchitectureName" json:"architecture_name"`
	TopFloor         string `xml:"ArchitectureStorys" json:"top_floor"` // 最高层数,以后可能会用到
	LowFloor         string `xml:"ArchitectureBegin" json:"low_floor"`  // 最低的层数
}

// 获取每个区域楼层信息
// 西区学生宿舍,东区学生宿舍,元宝山学生宿舍,南湖学生宿舍,国际园区 分别是1 2 3 4 5.
func GetArchitectures(area string) (*ArchiList, error) {
	var areaId string
	switch area {
	case "1":
		areaId = "0001"
	case "2":
		areaId = "0002"
	case "3":
		areaId = "0003"
	case "4":
		areaId = "0004"
	case "5":
		areaId = "0006"
	}

	var result ArchiList
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getArchitectureInfo?Area_ID=" + areaId

	err := getData(url, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 获取寝室信息
func GetRoomId(architectureId string, floorId string) (*DormList, error) {
	var result DormList
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getRoomInfo?Architecture_ID=" + architectureId + "&Floor=" + floorId

	err := getData(url, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetMeterInfo(roomId string) (*model.MeterInfo, error) {
	var result model.MeterInfo
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getMeterInfo?Room_ID=" + roomId

	err := getData(url, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

//  获取电表数据
func GetMeterInfoAPI(area, architecture, floor, dormName string) (*model.MeterInfo, error) {
	var (
		architectureId string
		dormInfo       DormInfo
		meterInfo      *model.MeterInfo
	)

	// 获取楼栋ID
	archiList, err := GetArchitectures(area)
	if err != nil {
		return nil, err
	}

	for _, v := range archiList.Architectures {
		if v.ArchitectureName == architecture {
			architectureId = v.ArchitectureID
			break
		}
	}

	// 获取宿舍ID
	dormList, err := GetRoomId(architectureId, floor)
	if err != nil {
		return nil, err
	}

	for _, v := range dormList.Dorms {
		if v.DormName == dormName {
			dormInfo.DormId = v.DormId
			dormInfo.DormName = dormName
			break
		}
	}

	// 获取电表ID
	meterInfo, err = GetMeterInfo(dormInfo.DormId)
	if err != nil {
		return nil, err
	}

	meterInfo.DormName = dormInfo.DormName

	return meterInfo, nil
}

// 有两个请求 第一个请求获取电费和时间，第二个请求获取昨日用量
func GetElectricCharge(meterId string, meterType string) (*model.ElectricCharge, error) {
	var charge model.ElectricCharge
	charge.Type = meterType
	charge.Id = meterId
	// 第二个URL的 query要求以“2020/5/30”格式输入昨天的日期
	// firstUrl 是第一个请求的URL,secondUrl 是第二个
	now := time.Now().AddDate(0, 0, -1)
	query := now.Format("2006/01/02")

	firstUrl := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getReserveHKAM?AmMeter_ID=" + meterId

	secondUrl := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getMeterDayValue?AmMeter_ID=" + meterId + "&startDate=" + query + "&endDate=" + query

	// 构造第一个请求
	err1 := getData(firstUrl, &charge)

	// 构造第二个请求
	err2 := getData(secondUrl, &charge)

	if err1 != nil || err2 != nil {
		log.Error("getCharge function error")
		cache, err := model.GetElectricity(meterId)
		if err != nil {
			return nil, err
		}
		return cache, nil
	}

	if err := model.AddElectricity(charge); err != nil {
		log.Error("AddElectricity function error" + err.Error())
	}

	return &charge, nil
}

// 爬取分析数据
func getData(url string, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(body, data)
	if err != nil {
		return err
	}
	return nil
}
