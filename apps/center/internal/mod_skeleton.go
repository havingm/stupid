package internal

import (
	"github.com/stupid/chanrpc"
	"github.com/stupid/module"
)

var (
	Skeleton = NewSkeleton()
	ChanRPC  = Skeleton.ChanRPCServer
)

// NewSkeleton todo 参数从配置读取
func NewSkeleton() *module.Skeleton {
	skn := &module.Skeleton{
		GoLen:              2048,
		TimerDispatcherLen: 2048,
		AsynCallLen:        2048,
		ChanRPCServer:      chanrpc.NewServer(2048),
	}
	skn.Init()
	return skn
}

type ModSkeleton struct {
	*module.Skeleton
}

func (m *ModSkeleton) OnInit() {
	m.Skeleton = Skeleton

}

func (m *ModSkeleton) OnDestroy() {

}
