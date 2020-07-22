package service

import (
	"encoding/xml"
	"errors"
	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/model"
	"io/ioutil"
	"net/http"
	"strconv"
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
	TopFloor         string `xml:"ArchitectureStorys" json:"top_floor"`  // 最高层数,以后可能会用到
	BeginFloor       string `xml:"ArchitectureBegin" json:"begin_floor"` // 最低的层数
}

// 获取每个区域楼层信息
// 西区学生宿舍,东区学生宿舍,元宝山学生宿舍,南湖学生宿舍,国际园区 分别是1 2 3 4 5.
func GetArchitectures(Area string) (ArchiList, error) {
	var areaId string
	switch Area {
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
	if areaId == "" {
		return ArchiList{}, errors.New("wrong query parameter")
	}

	var result ArchiList
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getArchitectureInfo?Area_ID=" + areaId

	resp, err := http.Get(url)
	if err != nil {
		return ArchiList{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ArchiList{}, err
	}

	err = xml.Unmarshal(body, &result)
	if err != nil {
		return ArchiList{}, err
	}
	return result, nil
}

// 获取寝室信息
func GetRoomId(ArchitectureId string, FloorId string) (DormList, error) {
	var result DormList
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getRoomInfo?Architecture_ID=" + ArchitectureId + "&Floor=" + FloorId
	resp, err := http.Get(url)
	if err != nil {
		return DormList{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DormList{}, err
	}

	err = xml.Unmarshal(body, &result)
	if err != nil {
		return DormList{}, err
	}
	return result, nil
}

func GetMeterInfo(RoomId string) (model.MeterInfo, error) {
	var result model.MeterInfo
	url := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getMeterInfo?Room_ID=" + RoomId
	resp, err := http.Get(url)
	if err != nil {
		return model.MeterInfo{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.MeterInfo{}, err
	}

	err = xml.Unmarshal(body, &result)
	if err != nil {
		return model.MeterInfo{}, err
	}
	return result, nil
}

//  获取电表数据
func GetMeterInfoAPI(Area, Architecture, Floor, DormName string) (model.MeterInfo, error) {
	var (
		architectureId string
		dormInfo       DormInfo
		meterInfo      model.MeterInfo
	)

	// 获取楼栋ID
	archiList, err := GetArchitectures(Area)
	if err != nil {
		return model.MeterInfo{}, err
	}

	for _, v := range archiList.Architectures {
		if v.ArchitectureName == Architecture {
			architectureId = v.ArchitectureID
			break
		}
	}

	// 获取宿舍ID
	dormList, err := GetRoomId(architectureId, Floor)
	if err != nil {
		return model.MeterInfo{}, err
	}

	for _, v := range dormList.Dorms {
		if v.DormName == DormName {
			dormInfo.DormId = v.DormId
			dormInfo.DormName = DormName
			break
		}
	}

	// 获取电表ID
	meterInfo, err = GetMeterInfo(dormInfo.DormId)
	if err != nil {
		return model.MeterInfo{}, err
	}

	meterInfo.DormName = dormInfo.DormName

	return meterInfo, nil
}

// 有两个请求 第一个请求获取电费和时间，第二个请求获取昨日用量
func GetElectricCharge(RoomId string, meterType string) (model.ElectricCharge, error) {
	var charge model.ElectricCharge
	charge.Type = meterType
	charge.Id = RoomId
	// 第二个URL的 query要求以“2020/5/30”格式输入昨天的日期
	// firstUrl 是第一个请求的URL,secondUrl 是第二个
	now := time.Now().AddDate(0, 0, -1)
	year, month, day := now.Date()

	//	query := now.Format("2020/7/17") 这种写法没得到预期的输入
	query := strconv.Itoa(year) + "%2F" + strconv.Itoa(int(month)) + "%2F" + strconv.Itoa(day)

	firstUrl := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getReserveHKAM?AmMeter_ID=" + RoomId

	secondUrl := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/" +
		"getMeterDayValue?AmMeter_ID=" + RoomId + "&startDate=" + query + "&endDate=" + query

	// 构造第一个请求
	resp, err := http.Get(firstUrl)
	if err != nil {
		// 请求失败就从数据库中获取数据,下面同理
		cache, err := model.GetElectricity(RoomId)
		if err != nil {
			return charge, err
		}
		return cache, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cache, err := model.GetElectricity(RoomId)
		if err != nil {
			return charge, err
		}
		return cache, nil
	}

	err = xml.Unmarshal(body, &charge)
	if err != nil {
		cache, err := model.GetElectricity(RoomId)
		if err != nil {
			return charge, err
		}
		return cache, nil
	}

	// 构造第二个请求
	res, err := http.Get(secondUrl)
	if err != nil {
		cache, err := model.GetElectricity(RoomId)
		if err != nil {
			return charge, err
		}
		return cache, nil
	}

	defer res.Body.Close()
	newBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		cache, err := model.GetElectricity(RoomId)
		if err != nil {
			return charge, err
		}
		return cache, nil
	}
	err = xml.Unmarshal(newBody, &charge)
	if err != nil {
		cache, err := model.GetElectricity(RoomId)
		if err != nil {
			return charge, err
		}
		return cache, nil
	}

	if err = model.AddElectricity(charge); err != nil {
		log.Error("AddElectricity function error" + err.Error())
	}

	return charge, nil
}
