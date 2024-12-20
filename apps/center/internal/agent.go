package internal

import "net"

type Agent interface {
	WriteMsg(msg interface{})
	WriteBytes(data []byte)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	GetData() interface{}
	SetData(data interface{})
}
