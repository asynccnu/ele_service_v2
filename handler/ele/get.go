package ele

import (
	"github.com/asynccnu/ele_service_v2/handler"
	"github.com/asynccnu/ele_service_v2/log"
	"github.com/asynccnu/ele_service_v2/service"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	id:=c.Query("AmMeter_ID")
	result,err:=service.GetElectricCharge(id)
	if err!=nil{
		log.Error("GetElectricityCharge function error")
		handler.SendError(c,err,nil,err.Error())
		return
	}
	handler.SendResponse(c,nil,result)
}
