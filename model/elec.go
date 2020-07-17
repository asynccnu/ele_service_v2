package model

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	MongoDb  = "electricity"
	ElecCol  = "charge"
	MeterCol = "meter"
)

// 电费储存结构
type ElectricCharge struct {
	Id           string           `bson:"meter_id"`                 // 电表ID
	RemainPower  string           `xml:"remainPower"  bson:"power"` // 剩余电量 单位：度
	ReadTime     string           `xml:"readTime"  bson:"time"`     // 最近一次抄表时间
	ElectricInfo ElectricInfoList `xml:"dayValueInfoList>DayValueInfo"`
}

//用List是因为这个接口返回一段时间内的数据,但是我把日期限制在昨天,就只返回昨天的数据
// 电费储存结构
type ElectricInfoList struct {
	YesterdayElecUse string `xml:"dayValue" bson:"value"`     // 昨日用电量
	YesterdayFee     string `xml:"dayUseMeony"  bson:"money"` // 昨日电费
}

// 电表号储存结构
type MeterInfo struct {
	DormName string `bson:"dorm_name" xml:"-"`
	MeterId  string `bson:"meter_id" xml:"meterList>MeterInfo>meterId"`
}

//获取mongodb中的寝室电表信息,如果没有则返回错误
func GetMongoMeterInfo(DormName string) (MeterInfo, error) {
	collection := DB.Self.Database(MongoDb).Collection(MeterCol)
	var result MeterInfo
	cur, err := collection.Find(context.TODO(), bson.M{"dorm_name": DormName})
	if err != nil {
		return MeterInfo{}, nil
	}
	for cur.Next(context.TODO()) {
		err := cur.Decode(&result)
		if err != nil {
			return MeterInfo{}, err
		}
	}

	return result, nil
}

// 查看mongodb中是否有电表信息
func HaveMeterInfo(dormName string) (bool, error) {
	collection := DB.Self.Database(MongoDb).Collection(MeterCol)

	count, err := collection.CountDocuments(context.TODO(), bson.M{"dorm_name": dormName})

	if err != nil {
		return false, err
	}

	return count != 0, nil
}

// 添加电表信息
func AddMeterInfo(meter MeterInfo) error {
	collection := DB.Self.Database(MongoDb).Collection(MeterCol)
	var err error

	if haveDoc, err := HaveMeterInfo(meter.DormName); !haveDoc {
		if err != nil {
			return err
		}
		_, err = collection.InsertOne(context.TODO(), meter)
	}

	return err
}

// 查看mongodb中是否电费记录
func HaveElectricity(id string) (bool, error) {
	collection := DB.Self.Database(MongoDb).Collection(ElecCol)

	count, err := collection.CountDocuments(context.TODO(), bson.M{"meter_id": id})

	if err != nil {
		return false, err
	}

	return count != 0, nil
}

// 插入或更新电费记录
func AddElectricity(rec ElectricCharge) error {
	collection := DB.Self.Database(MongoDb).Collection(ElecCol)
	var err error

	// 有记录则为替换，无记录就插入
	if haveDoc, _ := HaveElectricity(rec.Id); haveDoc {
		_, err = collection.ReplaceOne(context.TODO(), bson.M{"meter_id": rec.Id}, rec)

	} else {
		_, err = collection.InsertOne(context.TODO(), rec)
	}

	return err
}

// 从数据库中获取电费记录
func GetElectricity(id string) (ElectricCharge, error) {
	var result ElectricCharge

	collection := DB.Self.Database(MongoDb).Collection(ElecCol)

	if haveDoc, _ := HaveElectricity(id); !haveDoc {
		return result, errors.New("get electricity charge failed and can't find document" +
			" from MongoDB")
	}

	cur, err := collection.Find(context.TODO(), bson.M{"meter_id": id})
	if err != nil {
		return result, err
	}

	for cur.Next(context.TODO()) {
		err := cur.Decode(&result)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}
