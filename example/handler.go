package example

/*
import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"go-admin/app/server/SmatLocker"
	"go-admin/app/server/WearDevice"
	"go-admin/pkg/logger"
	"io"
	"net"
	"runtime/debug"
	"strings"
)

const (
	App_WearDevice = 0x01
	App_SmatLocker = 0x02
)

func checkStart(rw io.Reader, Start []byte) (err error) {
	i := 0
	var b [1]byte
	n := 0

	if Start != nil && len(Start) > 0 {
		for {
			n, err = rw.Read(b[:])
			//err = binary.Read(r, binary.LittleEndian, &b)
			if err == io.EOF {
				return
			}
			if err != nil {
				return fmt.Errorf("Read Head %v\n", err)
			}
			if n == 0 {
				return fmt.Errorf("Read EOF \n")
			}
			if b[0] == Start[i] {
				i++
				if i == len(Start) {
					break
				}
			} else {
				i = 0
			}
		}
	}
	return
}
func EquipReceive(conn net.Conn) {
	defer func() {
		//<-ts.MaxConnChan
		if e := recover(); e != nil {
			err := fmt.Errorf("Panicing in EquipReceive %v\n", e)
			logger.TcpLogger.Errorf(err.Error())

			debug.PrintStack()
			return
		}
		conn.Close()

	}()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		err := checkStart(rw, []byte{0x55, 0xAA})
		if err == io.EOF {
			return
		}
		if err != nil {
			logger.TcpLogger.Error(err)
			return
		}
		var Len uint16
		appCode, _ := rw.ReadByte()
		funCode, _ := rw.ReadByte()
		switch appCode {
		case App_WearDevice:
			err = binary.Read(rw, binary.BigEndian, &Len)
			if err != nil {
				logger.TcpLogger.Error(err)
				return
			}
		case App_SmatLocker:
			err = binary.Read(rw, binary.LittleEndian, &Len)
			if err != nil {
				logger.TcpLogger.Error(err)
				return
			}
		default:
			logger.TcpLogger.Error(fmt.Errorf("AppCode err :0x%02X", appCode))
			continue
		}

		if Len > 1024 {
			logger.TcpLogger.Warningf("read package len is too big > 1024.")
		}
		buffer := make([]byte, Len)

		rw.Read(buffer)
		logger.TcpLogger.Debugf("received [%d] from[%v]:%s", len(buffer), conn.RemoteAddr(), strings.ToUpper(hex.EncodeToString(buffer)))
		respbuf := []byte{0x55, 0xAA, appCode, funCode, 0x00, 0x01}
		var resp []byte
		switch appCode {
		case App_WearDevice:
			resp, err = WearDevice.ParseData(funCode, buffer)
		case App_SmatLocker:
			resp, err = SmatLocker.ParseData(funCode, buffer)
		}

		if err != nil {
			logger.TcpLogger.Error(err)
			respbuf = append(respbuf, 0x01)
		} else {
			respbuf = append(respbuf, 0x00)
			respbuf = append(respbuf, resp...)
		}
		rw.Write(respbuf)
		rw.Flush()
	}
}
*/