package internal

import (
	"github.com/stupid/log"
)

func init() {

	Skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	Skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(Agent)
	log.Debug("rpcNewAgent agent:%v", a)
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(Agent)
	aData := a.GetData()
	log.Debug("rpcCloseAgent,userData:%v", aData)
}
