package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Requet *gin.Context
}

type ResponseData struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func responseData(code int, msg string, data interface{}, response *Response)  {
	response.Requet.JSON(http.StatusOK, ResponseData{code, msg, data})
}

func (response *Response) Success ()  {
	responseData(200, "OK", map[string]interface{}{}, response)
}

func (response *Response) SuccessWithData(data interface{})  {
	responseData(200, "OK", data, response)
}

func (response *Response) SuccessWithMSG(msg string)  {
	responseData(200, msg, map[string]interface{}{}, response)
}

func (response *Response) SuccessWithDetail(code int, msg string, data interface{})  {
	responseData(code, msg, data, response)
}

func (response *Response) Error()  {
	responseData(500, "Error", map[string]interface{}{}, response)
}

func (response *Response) ErrorWithMSG(msg string)  {
	responseData(500, msg, map[string]interface{}{}, response)
}

func (response *Response) ErrorWithData(data interface{})  {
	responseData(500, "Error", data, response)
}

func (response *Response) ErrorWithDetail(code int, msg string, data interface{})  {
	responseData(code, msg, data, response)
}