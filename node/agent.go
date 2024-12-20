package node

import (
	"github.com/stupid/log"
	"github.com/stupid/network"
	"net"
	"reflect"
)

type agent struct {
	conn network.Conn
	node *Node
	data interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.node.Processor != nil {
			msg, e := a.node.Processor.Unmarshal(data)
			if e != nil {
				log.Debug("unmarshal message error: %v", e)
				break
			}
			e = a.node.Processor.Route(msg, a)
			if e != nil {
				log.Debug("route message error: %v", e)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.node.Skeleton.ChanRPCServer != nil {
		err := a.node.Skeleton.ChanRPCServer.Call0("CloseAgent", a)
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
	if a.node.Processor != nil {
		data, err := a.node.Processor.Marshal(msg)
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
