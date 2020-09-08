package ele

import (
	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/pkg/errno"

	"github.com/gin-gonic/gin"
)

type BuildingItem struct {
	Name  string `json:"name"`  // 楼栋名，用于用户展示
	Alias string `json:"alias"` // 别名简称，用于请求
}

// 西区
var WestAreaBuildings = []BuildingItem{
	{Name: "西区1栋", Alias: "西1"},
	{Name: "西区2栋", Alias: "西2"},
	{Name: "西区3栋", Alias: "西3"},
	{Name: "西区4栋", Alias: "西4"},
	{Name: "西区5栋", Alias: "西5"},
	{Name: "西区6栋", Alias: "西6"},
	{Name: "西区7栋", Alias: "西7"},
	{Name: "西区8栋", Alias: "西8"},
}

// 东区
var EastAreaBuildings = []BuildingItem{
	{Name: "东区1栋", Alias: "东1"},
	{Name: "东区2栋", Alias: "东2"},
	{Name: "东区3栋", Alias: "东3"},
	{Name: "东区4栋", Alias: "东4"},
	{Name: "东区5栋", Alias: "东5"},
	{Name: "东区6栋", Alias: "东6"},
	{Name: "东区7栋", Alias: "东7"},
	{Name: "东区8栋", Alias: "东8"},
	{Name: "东区9栋", Alias: "东9"},
	{Name: "东区10栋", Alias: "东10"},
	{Name: "东区11栋", Alias: "东11"},
	{Name: "东区12栋", Alias: "东12"},
	{Name: "东区13栋东", Alias: "东13-东"},
	{Name: "东区13栋西", Alias: "东13-西"},
	{Name: "东区14栋", Alias: "东14"},
	{Name: "东区15栋东", Alias: "东15-东"},
	{Name: "东区15栋西", Alias: "东15-西"},
	{Name: "东区16栋", Alias: "东16"},
	{Name: "东区附1栋", Alias: "东附1"},
}

// 元宝山
var YuanBaoMountainBuildings = []BuildingItem{
	{Name: "元宝山1栋", Alias: "元1"},
	{Name: "元宝山2栋", Alias: "元2"},
	{Name: "元宝山3栋", Alias: "元3"},
	{Name: "元宝山4栋", Alias: "元4"},
	{Name: "元宝山5栋", Alias: "元5"},
}

// 南湖
var SouthLakeBuildings = []BuildingItem{
	{Name: "南湖1栋", Alias: "南湖1"},
	{Name: "南湖2栋", Alias: "南湖2"},
	{Name: "南湖3栋", Alias: "南湖3"},
	{Name: "南湖4栋", Alias: "南湖4"},
	{Name: "南湖5栋", Alias: "南湖5"},
	{Name: "南湖6栋", Alias: "南湖6"},
	{Name: "南湖7栋", Alias: "南湖7"},
	{Name: "南湖8栋", Alias: "南湖8"},
	{Name: "南湖9栋", Alias: "南湖9"},
	{Name: "南湖10栋", Alias: "南湖10"},
	{Name: "南湖11栋", Alias: "南湖11"},
	{Name: "南湖12栋", Alias: "南湖12"},
	{Name: "南湖13栋", Alias: "南湖13"},
}

// 国交
var InternationalExchangeBuildings = []BuildingItem{
	{Name: "国交4栋", Alias: "国4"},
	{Name: "国交8栋", Alias: "国8"},
	{Name: "国交9栋", Alias: "国9"},
}

// 获取楼栋
// area: 西区/东区/元宝山/南湖/国交
func GetBuildings(c *gin.Context) {
	area := c.DefaultQuery("area", "")
	if area == "" {
		handler.SendBadRequest(c, errno.ErrQuery, nil, "area is required")
		return
	}

	var list []BuildingItem

	switch area {
	case "西区":
		list = WestAreaBuildings
	case "东区":
		list = EastAreaBuildings
	case "元宝山":
		list = YuanBaoMountainBuildings
	case "南湖":
		list = SouthLakeBuildings
	case "国交":
		list = InternationalExchangeBuildings
	default:
		handler.SendBadRequest(c, errno.ErrQuery, nil, "the area does not exist.")
		return
	}

	handler.SendResponse(c, nil, list)
}
