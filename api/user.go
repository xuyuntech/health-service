package api

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *Api) queryAllUser(c *gin.Context) {
	queryString := c.Param("query_string")
	if strings.Trim(queryString, " ") == "" {
		RespErr(c, errors.New("need query_string param"))
		return
	}
	payload, err := a.Fabric.QueryUser(queryString)
	if err != nil {
		RespErr(c, err)
		return
	}
	Resp(c, payload)
}

func (a *Api) initUser(c *gin.Context) {
	txid, err := a.Fabric.InitUser()
	if err != nil {
		RespErr(c, err)
		return
	}
	Resp(c, txid)
}
