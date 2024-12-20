package node

import (
	"github.com/stupid/chanrpc"
	"github.com/stupid/module"
	"github.com/stupid/network"
	"github.com/stupid/processor/protobuf"
)

const (
	kGoLen              = 2048
	kTimerDispatcherLen = 2048
	kAsynCallLen        = 2048
	kCallChanLen        = 2048
)

const (
	kMaxConnections  = 2048
	kPendingWriteNum = 2048
	kLenMsgLen       = 4
	kMaxMsgLen       = 204800
	kLittleEndian    = false
)

type Node struct {
	Skeleton  *module.Skeleton
	ServAddr  string
	Processor network.Processor
}

func NewNode(servAddr string) *Node {
	n := &Node{}
	n.ServAddr = servAddr
	return n
}

func (n *Node) OnInit() {
	rpcServer := chanrpc.NewServer(kCallChanLen)
	skn := &module.Skeleton{
		GoLen:              kGoLen,
		TimerDispatcherLen: kTimerDispatcherLen,
		AsynCallLen:        kAsynCallLen,
		ChanRPCServer:      rpcServer,
	}
	skn.Init()
	n.Skeleton = skn
	n.Processor = protobuf.NewProcessor()
}

func (n *Node) OnDestroy() {

}

func (n *Node) Run(closeSig chan bool) {
	go n.Skeleton.Run(closeSig)
	var server *network.TCPServer
	if n.ServAddr != "" {
		server = new(network.TCPServer)
		server.Addr = n.ServAddr
		server.MaxConnNum = kMaxConnections
		server.PendingWriteNum = kPendingWriteNum
		server.LenMsgLen = kLenMsgLen
		server.MaxMsgLen = kMaxMsgLen
		server.LittleEndian = kLittleEndian
		server.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, node: n}
			if n.Skeleton != nil {
				n.Skeleton.ChanRPCServer.Go("NewAgent", a)
			}
			return a
		}
	}
	if server != nil {
		server.Start()
	}
	<-closeSig
	if server != nil {
		server.Close()
	}
}
