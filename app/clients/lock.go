package clients

import (
	"bytes"
	"encoding/hex"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/util/gconv"
	"locksim/app/uitls"
	"path"
	"time"
)

var Lock = new(LockCmd)

type LockCmd struct {
	CmdType string
	CMD30   *ActiveControl
	CMD31   *SearchPassword
	CMD32   *SearchCache
	CMD33   *ReportEvent
	CMD34   *ReportImage
	CMD39   *IssuedOpen
	CMD40   *IssuedLock
}

// 设备激活平台控制功能 30
type ActiveControl struct {
	DeviceMac        string `p:"DeviceMac"  v:"required|length:12,12"`      //设备 Mac 地址，HEX 格式
	VerificationCode string `p:"VerificationCode"  v:"required|length:6,6"` //用户输入的验证码，ASCII 格式
}

// 查询平台开锁密码 31
type SearchPassword struct {
	DeviceID string `p:"DeviceID"  v:"required|length:12,12"` //设备编码，ASCII 格式
}

// 查询平台缓存指令 32
type SearchCache struct {
	DeviceID string `p:"DeviceID"  v:"required|length:12,12"` //设备编码，ASCII 格式
}

// 主动上报事件 33
type ReportEvent struct {
	DeviceID  string `p:"DeviceID"  v:"required|length:12,12"` //设备编码，ASCII 格式
	EventType uint8  `p:"EventType"  v:"required|between:0,4"` //事件类型
	EventCode uint8  `p:"EventCode"  v:"required|length:1,1"`  //事件编号
	UserType  uint8  `p:"UserType"  v:"required|between:0,9"`  //用户类型
	UserID    uint16 `p:"UserID"  v:"required|between:0,2048"` //用户编号，Word 格式 本地使用的范围: 0-1024 云用户使用的范围: 1024-2048
}

// 主动上报人脸图像 34
type ReportImage struct {
	Image     *ghttp.UploadFile //人脸图像文件
	DeviceID  string            `p:"DeviceID"  v:"required|length:12,12"` //设备编码，ASCII 格式
	ImageType uint8             //图像格式：0x00-bmp 0x01-jpeg 0x02-png 0x03-gif
	ImageLen  uint16            //图像数据长度
	ImageData []byte            //图像数据
}

// 下发开锁指令 39
type IssuedOpen struct {
	DeviceID  string `p:"DeviceID"  v:"required|length:12,12"` //设备编码，ASCII 格式
	Timestamp uint32 `p:"Timestamp"  v:"required"`             //时间戳，DWORD 格式，UTC 时间秒数
	Sign      string `p:"Sign"  v:"required|length:4,4"`       //签名数据最后 4 字节内容
}

// 下发锁死/解锁指令 40
type IssuedLock struct {
	DeviceID  string `p:"DeviceID"  v:"required|length:12,12"` //设备编码，ASCII 格式
	Timestamp uint32 `p:"Timestamp"  v:"required"`             //时间戳，DWORD 格式，UTC 时间秒数
	Disable   string `p:"Sign"  v:"required|between:0,2"`      //0x00–解锁 0x02-锁死
	Sign      string `p:"Sign"  v:"required|length:4,4"`       //签名数据最后4字节内容
}

// 向服务端发送指令
func (s *LockCmd) SendPacket() (string, error) {
	pkg := s.makePacket(s.CmdType)
	conn, err := gtcp.NewConn(config.GetString("address"), time.Second*2)
	if err != nil {
		logger.Line().Error(err)
		return "", err
	}
	logger.Info("发送:", hex.EncodeToString(pkg))
	if res, err := conn.SendRecvWithTimeout(pkg, -1, time.Second*3); err == nil {
		logger.Info("返回:", hex.EncodeToString(res))
		return string(res), nil
	} else {
		return "", err
	}
	return "", err
}

func (s *LockCmd) makePacket(msg string) []byte {
	// 同步1byte,起始1byte,应用码1byte
	b := []byte{0x55, 0xAA, 0x02}
	// 数据长度 2byte
	var mlen uint16
	switch msg {
	case "30":
		b = append(b, 0x30)
		MacString := s.CMD30.DeviceMac
		MacByte, _ := hex.DecodeString(MacString)
		vefByte := []byte(s.CMD30.VerificationCode)
		mlen = gconv.Uint16(len(MacByte) + len(vefByte))
		logger.Info("len30", mlen)
		b = append(b, uitls.ConverNumberToBytes(mlen)...)
		b = append(b, MacByte...)
		b = append(b, vefByte...)
	case "31":
		b = append(b, 0x31)
		did := []byte(s.CMD31.DeviceID)
		mlen = gconv.Uint16(len(did))
		logger.Info("len31", mlen)
		b = append(b, uitls.ConverNumberToBytes(mlen)...)
		b = append(b, did...)
	case "32":
		b = append(b, 0x32)
		did := []byte(s.CMD32.DeviceID)
		mlen = gconv.Uint16(len(did))
		logger.Info("len32", mlen)
		b = append(b, uitls.ConverNumberToBytes(mlen)...)
		b = append(b, did...)
	case "33":
		b = append(b, 0x33)
		did := []byte(s.CMD33.DeviceID)
		uid := uitls.ConverNumberToBytes(s.CMD33.UserID)
		mlen = gconv.Uint16(len(did) + 3 + len(uid))
		logger.Info("len33", mlen)
		b = append(b, uitls.ConverNumberToBytes(mlen)...)
		b = append(b, did...)
		b = append(b, uitls.ConverNumberToBytes(s.CMD33.EventType)...)
		b = append(b, uitls.ConverNumberToBytes(s.CMD33.EventCode)...)
		b = append(b, uitls.ConverNumberToBytes(s.CMD33.UserType)...)
		b = append(b, uid...)
	case "34":
		b = append(b, 0x34)
		f, err := s.CMD34.Image.Open()
		if err != nil {
			logger.Error(err.Error())
		}
		defer f.Close()
		buff := &bytes.Buffer{}
		buff.ReadFrom(f)

		s.CMD34.ImageData = buff.Bytes()
		s.CMD34.ImageLen = gconv.Uint16(len(buff.Bytes()))
		suffix := path.Ext(s.CMD34.Image.Filename)
		switch suffix {
		case ".bmp":
			s.CMD34.ImageType = 0
		case ".jpeg":
			s.CMD34.ImageType = 1
		case ".png":
			s.CMD34.ImageType = 2
		case ".gif":
			s.CMD34.ImageType = 3
		default:
			s.CMD34.ImageType = 4
		}
		did := []byte(s.CMD34.DeviceID)
		iml := uitls.ConverNumberToBytes(s.CMD34.ImageLen)
		mlen = gconv.Uint16(len(did) + 1 + len(iml) + len(s.CMD34.ImageData))
		logger.Info("len34", mlen)
		b = append(b, uitls.ConverNumberToBytes(mlen)...)
		b = append(b, did...)
		b = append(b, uitls.ConverNumberToBytes(s.CMD34.ImageType)...)
		b = append(b, iml...)
		b = append(b, s.CMD34.ImageData...)
	}

	crc := uitls.CRCCheck(b, len(b))
	b = append(b, uitls.ConverNumberToBytes(crc)...) //校验CRC
	return b
}
