package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

type Response struct {
	Requet *gin.Context
}

type ResponseData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func responseData(code int, msg string, data interface{}, response *Response) {
	response.Requet.JSON(http.StatusOK, ResponseData{code, msg, data})
}

func (response *Response) Success() {
	responseData(200, "OK", map[string]interface{}{}, response)
}

func (response *Response) SuccessWithData(data interface{}) {
	responseData(200, "OK", data, response)
}

func (response *Response) SuccessWithMSG(msg string) {
	responseData(200, msg, map[string]interface{}{}, response)
}

func (response *Response) SuccessWithDetail(code int, msg string, data interface{}) {
	responseData(code, msg, data, response)
}

func (response *Response) SuccessWithFile(file string) {

	var HttpContentType = map[string]string{
		".avi": "video/avi",
		".mp3": "   audio/mp3",
		".mp4": "video/mp4",
		".wmv": "   video/x-ms-wmv",
		".asf":  "video/x-ms-asf",
		".rm":   "application/vnd.rn-realmedia",
		".rmvb": "application/vnd.rn-realmedia-vbr",
		".mov":  "video/quicktime",
		".m4v":  "video/mp4",
		".flv":  "video/x-flv",
		".jpg":  "image/jpeg",
		".png":  "image/png",
	}

	fileNameWithSuffix := path.Base(file)
	fileType := path.Ext(fileNameWithSuffix)
	fileContentType, ok := HttpContentType[fileType]
	if !ok {
		response.ErrorWithMSG("头像加载失败")
		return
	}

	response.Requet.Header("Content-Type", fileContentType)
	response.Requet.File(file)
}

func (response *Response) Error() {
	responseData(500, "Error", map[string]interface{}{}, response)
}

func (response *Response) ErrorWithMSG(msg string) {
	responseData(500, msg, map[string]interface{}{}, response)
}

func (response *Response) ErrorWithData(data interface{}) {
	responseData(500, "Error", data, response)
}

func (response *Response) ErrorWithDetail(code int, msg string, data interface{}) {
	responseData(code, msg, data, response)
}
