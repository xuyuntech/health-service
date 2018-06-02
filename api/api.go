package api

import (
	"errors"
	"strings"
	"time"

	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xuyuntech/health-service/blockchain"
)

type Api struct {
	Fabric *blockchain.FabricSetup
}

func (a *Api) Run() error {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Requested-With", "X-Access-Token"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/arrangement", a.arrangement)
	r.GET("/initData", a.initData)
	r.GET("/query", a.query)
	r.GET("/createRegister", a.createRegister)
	r.POST("/updateRegister", a.updateRegister)
	return r.Run()
}

func (a *Api) arrangement(c *gin.Context) {
	form := &blockchain.ArrangementForm{}
	if err := c.BindJSON(form); err != nil {
		RespErr(c, fmt.Errorf("参数错误(%v)", err))
		return
	}
	payload, err := a.Fabric.Arrangement(form)
	if err != nil {
		RespErr(c, err)
		return
	}
	Resp(c, payload)
}
func (a *Api) initData(c *gin.Context) {
	payload, err := a.Fabric.InitData()
	if err != nil {
		RespErr(c, err)
		return
	}
	Resp(c, payload)
}

func (a *Api) query(c *gin.Context) {
	queryString := c.Query("query_string")
	if strings.Trim(queryString, " ") == "" {
		RespErr(c, errors.New("need query_string param"))
		return
	}
	payload, err := a.Fabric.Query(queryString)
	if err != nil {
		RespErr(c, err)
		return
	}
	Resp(c, payload)
}
func (a *Api) createRegister(c *gin.Context) {
	// query: userKey, arrangementKey
	userKey := c.Query("userKey")
	arrangementKey := c.Query("arrangementKey")

	payload, err := a.Fabric.CreateRegister(userKey, arrangementKey)
	if err != nil {
		RespErr(c, err)
		return
	}
	Resp(c, payload)
}
func (a *Api) updateRegister(c *gin.Context) {
	// 参数说明
	// 0 		1					2		[3			4			5			6				7]
	// userKey	registerHistoryKey 	state 	complained 	diagnose 	history 	familyHistory 	items [][]string
	// 用户		挂号记录				状态		主诉			诊断			病史			家族史			药品列表
	updateRegisterForm := &blockchain.UpdateRegisterForm{}
	if err := c.BindJSON(updateRegisterForm); err != nil {
		RespErr(c, fmt.Errorf("参数错误(%v)", err))
		return
	}
	payload, err := a.Fabric.UpdateRegister(updateRegisterForm)
	if err != nil {
		RespErr(c, err)
		return
	}
	Resp(c, payload)
}

func RespErr(c *gin.Context, err error, msg ...string) {
	results := map[string]interface{}{
		"status": 1,
		"err":    err.Error(),
	}
	if len(msg) >= 1 {
		results["msg"] = msg[0]
	}
	c.JSON(200, results)
}

func Resp(c *gin.Context, results interface{}) {
	c.JSON(200, map[string]interface{}{
		"status": 0,
		"data":   results,
	})
}
