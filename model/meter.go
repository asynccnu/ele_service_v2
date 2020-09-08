package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// 电表号储存结构
type MeterModel struct {
	Building string       `bson:"building"` // 楼栋，如东16
	Room     string       `bson:"room"`     // 寝室，如101
	Meters   []*MeterInfo `bson:"meters"`
}

type MeterInfo struct {
	Kind    string `bson:"kind"`     // 类型，light/air
	MeterID string `bson:"meter_id"` // 电表号
}

// 获取寝室电表信息
func GetMetersByBuildingAndRoom(building, room string) ([]*MeterInfo, error) {
	var result MeterModel

	err := DB.Self.Database(DBName).Collection(MeterCol).
		FindOne(context.TODO(), bson.M{"building": building, "room": room}).
		Decode(&result)

	if err != nil {
		return nil, err
	}

	return result.Meters, nil
}

// 获取某宿舍楼中所有宿舍号
func GetRoomsByBuildingName(building string) ([]string, error) {
	var result []string

	cursor, err := DB.Self.Database(DBName).Collection(MeterCol).
		Find(context.TODO(), bson.M{"building": building})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var doc MeterModel
	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		result = append(result, doc.Room)
	}

	return result, err
}
