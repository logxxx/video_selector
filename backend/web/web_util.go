package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/logxxx/video_selector/base"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

func handleHome(c *gin.Context) {
	filePath := "dist/index.html"
	if !IsMobile(c.Request) {
		filePath = "dist_pc/index.html"
	}
	fileRealPath := filepath.Join(*base.WorkPath, filePath)

	c.File(fileRealPath)
}

func handleDist(c *gin.Context) {
	file := c.Param("filepath")
	isMobile := IsMobile(c.Request)
	distFolder := "dist"
	if !isMobile {
		distFolder = "dist_pc"
	}
	fileRealPath := filepath.Join(*base.WorkPath, distFolder, file)

	c.File(fileRealPath)
}

// check if is mobile broswer
func IsMobile(r *http.Request) bool {

	//put headers in a map
	headers := make(map[string]string)

	//beego中获取headers
	//headers := this.Ctx.Request.Header
	//for key,item := range headers{
	//fmt.Printf("%s=%s\n", key, item[0])
	//}

	//net/http中获取headers
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			headers[k] = v[0]
			//fmt.Printf("%s=%s\n", k, v[0])
		}
	}
	var is_mobile = false
	via := strings.ToLower(headers["VIA"])
	accept := strings.ToUpper(headers["Accept"])
	HTTP_X_WAP_PROFILE := headers["X_WAP_PROFILE"]
	HTTP_PROFILE := headers["PROFILE"]
	HTTP_USER_AGENT := headers["User-Agent"]
	if via != "" && strings.Index(via, "wap") != -1 {
		is_mobile = true
	} else if accept != "" && strings.Index(accept, "VND.WAP.WML") != -1 {
		is_mobile = true
	} else if HTTP_X_WAP_PROFILE != "" || HTTP_PROFILE != "" {
		is_mobile = true
	} else if HTTP_USER_AGENT != "" {

		reg := regexp.MustCompile(`(?i:(blackberry|configuration\/cldc|hp |hp-|htc |htc_|htc-|iemobile|kindle|midp|mmp|motorola|mobile|nokia|opera mini|opera |Googlebot-Mobile|YahooSeeker\/M1A1-R2D2|android|iphone|ipod|mobi|palm|palmos|pocket|portalmmm|ppc;|smartphone|sonyericsson|sqh|spv|symbian|treo|up.browser|up.link|vodafone|windows ce|xda |xda_|MicroMessenger))`)

		fmt.Printf("%q\n", reg.FindAllString(HTTP_USER_AGENT, -1))

		if len(reg.FindAllString(HTTP_USER_AGENT, -1)) > 0 {
			is_mobile = true
		}

	}

	return is_mobile

}
