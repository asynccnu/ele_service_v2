package service

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
type ArchitectureList struct {
	List []*ArchitectureInfo `xml:"architectureInfoList>architectureInfo" json:"architectures"`
}

// 楼栋信息
type ArchitectureInfo struct {
	ArchitectureID   string `xml:"ArchitectureID" json:"architecture_id"`
	ArchitectureName string `xml:"ArchitectureName" json:"architecture_name"`
	TopFloor         string `xml:"ArchitectureStorys" json:"top_floor"` // 最高层数,以后可能会用到
	LowFloor         string `xml:"ArchitectureBegin" json:"low_floor"`  // 最低的层数
}
