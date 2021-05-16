package api

import (
	"bytes"
	"encoding/base64"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"locksim/app/clients"
	"locksim/app/response"
)

var Lock = lockApi{}

type lockApi struct{}

func (a *lockApi) Send(r *ghttp.Request) {
	fun := r.Get("fun")
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

func (a *lockApi) Upload(r *ghttp.Request) {
	file := r.GetUploadFile("file")
	ctp := file.Header["Content-Type"][0]

	f, err := file.Open()
	if err != nil {
		r.Response.WriteExit(err)
	}
	defer f.Close()
	buff := &bytes.Buffer{}
	buff.ReadFrom(f)

	res := base64.StdEncoding.EncodeToString(buff.Bytes())
	r.Response.WriteExit("data:" + ctp + ";base64," + res)
}
