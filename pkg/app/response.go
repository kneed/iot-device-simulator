package app

import (
"github.com/gin-gonic/gin"
	"net/http"
)

type Gin struct {
	Ctx *gin.Context
}

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResponsePageBody struct {
	Total       int64       `json:"total"`
	TotalPage   int         `json:"totalPage"`
	PageSize    int         `json:"pageSize"`
	CurrentPage int         `json:"currentPage"`
	List        interface{} `json:"list"`
}

// Response settings gin.JSON
func (g *Gin) Response(errCode string, data interface{}) {
	g.Ctx.JSON(http.StatusOK, Response{
		Code: errCode,
		Msg:  GetMsg(errCode),
		Data: data,
	})
	return
}
