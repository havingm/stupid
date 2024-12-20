package internal

import (
	"fmt"
)

type ModCenter struct {
	*Center
}

func (m *ModCenter) OnInit() {
	m.Center = &Center{
		ServAddr:  fmt.Sprintf(":%d", 8000),
		Processor: Processor,
		ChanRPC:   ChanRPC,
	}
}
