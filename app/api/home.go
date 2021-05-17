package api

import (
	"github.com/gogf/gf/net/ghttp"
)

var Home = homeApi{}

type homeApi struct{}

// @summary 应用程序主页
// @description 只显示模板内容
// @tags    模拟锁端测试
// @produce html
// @router  /index [GET]
// @success 200 {string} string "执行结果"
func (a *homeApi) Index(r *ghttp.Request) {
	r.Response.WriteTpl("home/index.html")
}
