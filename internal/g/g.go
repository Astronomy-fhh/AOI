package g

import "time"

var (
	MapX                          = 1024
	MapY                          = 512
	MaxSize                       = 16
	MaxSplit                      = 4

	PlayerRadius          float64 = 10
	PlayerBaseSpeed               = 1e-7 // px/ns
	PlayerSpeedRandFactor         = 3
	PlayerMoveInternal            = time.Millisecond * 10
	TestMaxEnterPlayer            = 100

	ViewBackgroundColor           = "#FFFFFF"
	ViewTitle                     = "AOI"
	ViewScopeLineColor            = "#ACADA4"
	ViewScopePlayerColor          = "#DA49D3"
)
