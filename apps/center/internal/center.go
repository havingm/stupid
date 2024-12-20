package internal

import (
	"github.com/stupid/chanrpc"
	"github.com/stupid/log"
	"github.com/stupid/network"
	"net"
	"reflect"
)

const (
	kMaxConnections  = 1024
	kPendingWriteNum = 1024
	kLenMsgLen       = 4
	kMaxMsgLen       = 40960
	kLittleEndian    = false
)

type Center struct {
	ServAddr  string
	Processor network.Processor
	ChanRPC   *chanrpc.Server
}

func (c *Center) Run(closeSig chan bool) {
	var tcpServer *network.TCPServer
	if c.ServAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = c.ServAddr
		tcpServer.MaxConnNum = kMaxConnections
		tcpServer.PendingWriteNum = kPendingWriteNum
		tcpServer.LenMsgLen = kLenMsgLen
		tcpServer.MaxMsgLen = kMaxMsgLen
		tcpServer.LittleEndian = kLittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, center: c}
			if c.ChanRPC != nil {
				c.ChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if tcpServer != nil {
		tcpServer.Close()
	}

}

func (c *Center) OnDestroy() {}

type agent struct {
	conn   network.Conn
	center *Center
	data   interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.center.Processor != nil {
			msg, e := a.center.Processor.Unmarshal(data)
			if e != nil {
				log.Debug("unmarshal message error: %v", e)
				break
			}
			e = a.center.Processor.Route(msg, a)
			if e != nil {
				log.Debug("route message error: %v", e)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.center.ChanRPC != nil {
		err := a.center.ChanRPC.Call0("CloseAgent", a)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteBytes(data []byte) {
	err := a.conn.WriteMsg(data)
	if err != nil {
		log.Error("write bytes error: %v", err)
	}
}

func (a *agent) WriteMsg(msg interface{}) {
	if a.center.Processor != nil {
		data, err := a.center.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) GetData() interface{} {
	return a.data
}

func (a *agent) SetData(data interface{}) {
	a.data = data
}
