package api

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"locksim/app/clients"
	"locksim/app/response"
)

var Lock = lockApi{}

type lockApi struct{}

// @summary 执行报文发送
// @tags    锁端向应用程序发送报文
// @produce json
// @Param fun query int false "功能码" minimum(30) maximum(44)
// @router  /lock/send [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *lockApi) Send(r *ghttp.Request) {
	fun := r.Get("fun")
	if fun == nil {
		response.JsonExit(r, 1, "miss argument fun")
	}
	clock := new(clients.LockCmd)
	clock.CmdType = gconv.String(fun)
	var err error
	switch fun {
	case "30":
		err = r.Parse(&clock.CMD30)
	case "31":
		err = r.Parse(&clock.CMD31)
	case "32":
		err = r.Parse(&clock.CMD32)
	case "33":
		err = r.Parse(&clock.CMD33)
	case "34":
		err = r.Parse(&clock.CMD34)
		file := r.GetUploadFile("file")
		if file == nil {
			response.JsonExit(r, 1, "miss file")
		}
		clock.CMD34.Image = file
	default:
		response.JsonExit(r, 1, "指令未开放")
	}

	// 数据校验
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}

	// 执行发送
	res, err := clock.SendPacket()
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "success", res)
	}
}
