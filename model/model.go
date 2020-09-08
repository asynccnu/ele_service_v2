package model

import "go.mongodb.org/mongo-driver/mongo"

const (
	MongoDB  = "electricity"
	MeterCol = "meter"
	// ElecCol  = "ele"
)

var (
	MeterCollection *mongo.Collection
)

// 电费信息
type ElectricityInfo struct {
	Kind        string             `json:"kind"`         // 电费的类型
	RemainPower string             `json:"remain_power"` // 剩余电量 单位：度
	ReadTime    string             `json:"read_time"`    // 最近一次抄表时间
	Consumption EleConsumptionItem `json:"consumption"`  // 昨日消费
}

// 电费结构
type EleConsumptionItem struct {
	Usage  string `json:"usage"`  // 用电量
	Charge string `json:"charge"` // 电费
}
