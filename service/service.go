package service

import (
	"encoding/xml"
	"github.com/asynccnu/ele_service_v2/model"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//有两个请求 第一个请求获取电费和时间，第二个请求获取昨日用量
func GetElectricCharge(RoomId string) (model.ElectricCharge,error) {
	var charge model.ElectricCharge
	//第二个URL的 query要求以“2020/5/30”格式输入昨天的日期
	now := time.Now().AddDate(0,0,-1)
	year, month, day := now.Date()
	query := strconv.Itoa(year) + "%2F" + strconv.Itoa(int(month)) + "%2F" + strconv.Itoa(day)
	firUrl := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/getReserveHKAM?AmMeter_ID=" + RoomId
	secUrl := "http://jnb.ccnu.edu.cn/icbs/PurchaseWebService.asmx/getMeterDayValue?AmMeter_ID=" + RoomId + "&startDate=" + query + "&endDate=" + query
	//        构造第一个请求
	resp, err := http.Get(firUrl)
	if err!=nil{
		cache,err:=GetElectricCharge(RoomId)
		if err!=nil{
			return charge,err
		}
		return cache,nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		cache,err:=GetElectricCharge(RoomId)
		if err!=nil{
			return charge,err
		}
		return cache,nil
	}

	xml.Unmarshal(body, &charge)
	//        构造第二个请求
	res, err := http.Get(secUrl)
	if err!=nil{
		cache,err:=model.GetElectricity(RoomId)
		if err!=nil{
			return charge,err
		}
		return cache,nil
	}

	defer res.Body.Close()
	newBody, err := ioutil.ReadAll(res.Body)
	if err!=nil{
		cache,err:=model.GetElectricity(RoomId)
		if err!=nil{
			return charge,err
		}
		return cache,nil
	}
	charge.Id=RoomId
	xml.Unmarshal(newBody, &charge)
	model.AddEelctricity(charge)
	return charge,nil
}
