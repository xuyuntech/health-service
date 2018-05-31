package api

import (
	"github.com/gin-gonic/gin"
)

func (a *Api) queryAllHospitals(c *gin.Context) {
	payload, err := a.Fabric.QueryAllHospital()
	if err != nil {
		RespErr(c, err)
		return
	}
	Resp(c, payload)
}

func (a *Api) initHospital(c *gin.Context) {
	txid, err := a.Fabric.InitHospital()
	if err != nil {
		RespErr(c, err)
		return
	}
	Resp(c, txid)
}
