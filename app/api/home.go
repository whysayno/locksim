package api

import (
	"github.com/gogf/gf/net/ghttp"
)

var Home = homeApi{}

type homeApi struct{}

func (a *homeApi) Index(r *ghttp.Request) {
	r.Response.WriteTpl("home/index.html")
}
