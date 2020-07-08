package model

import ("context"
	"go.mongodb.org/mongo-driver/bson"
)
const (
	MongoDb = "ccnubox"
	ElecCol="electricity"
)
//电费储存结构
type ElectricCharge struct {
	Id string					`bson:"id"`
	RemainPower  string           `xml:"remainPower"  bson:"power"`  //剩余电量 单位：度
	ReadTime     string           `xml:"readTime"  bson:"time"`    //最近一次抄表时间
	DayValueInfo DayValueInfoList `xml:"dayValueInfoList>DayValueInfo"`
}
//电费储存结构
type DayValueInfoList struct {
	DayValue    string `xml:"dayValue" bson:"value"`    //昨日用电量
	DayUseMoney string `xml:"dayUseMeony"  bson:"money"` //昨日电费
}

// 查看monggodb中是否电费记录
func HaveElectricity(id string) (bool, error) {
	collection := DB.Self.Database(MongoDb).Collection(ElecCol)

	count, err := collection.CountDocuments(context.TODO(), bson.M{"id": id})

	if err != nil {
		return false, err
	} else if count == 0 {
		return false, nil
	}

	return true, nil
}

// 插入或更新电费记录
func AddEelctricity(rec ElectricCharge) error {
	collection := DB.Self.Database(MongoDb).Collection(ElecCol)
	var err error
	// 有记录则为替换，无记录就插入
	if haveDoc, _ := HaveElectricity(rec.Id); haveDoc {
		_, err = collection.ReplaceOne(context.TODO(), bson.M{"id": rec.Id}, rec)
	} else {
		_, err = collection.InsertOne(context.TODO(), rec)
	}

	if err != nil {
		return err
	}

	return nil
}

// 从数据库中获取电费记录
func GetElectricity(id string) ( ElectricCharge,error) {
	var result ElectricCharge

	collection := DB.Self.Database(MongoDb).Collection(ElecCol)

	cur, err := collection.Find(context.TODO(), bson.M{"id": id})
	if err != nil {
		return result,err
	}

	for cur.Next(context.TODO()) {
		err := cur.Decode(&result)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}