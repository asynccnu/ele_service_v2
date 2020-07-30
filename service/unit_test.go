package service

import (
	"github.com/asynccnu/ele_service_v2/model"
	"testing"
)

func TestGetElectricCharge(t *testing.T) {
	model.DB.Init()
	defer model.DB.Close()
	_,err:=GetElectricCharge("1001.001102.1")
	if err!=nil {
		t.Error(err)
	}

}