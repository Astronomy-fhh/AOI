package g

import "time"

var (
	MapX                                 = 1024
	MapY                                 = 512
	MaxSize                              = 16
	MaxSplit                             = 1
	PlayerRadius                 float64 = 10
	PlayerBaseSpeed                      = 1e-8 // px/ns
	PlayerSpeedRandFactor                = 5
	PlayerMoveInternal                   = time.Millisecond * 10
	PlayerMove                           = true
	TestMaxEnterPlayer                   = 1
	ViewTitle                            = "AOI"
	ViewBackgroundColor                  = "#FFFFFF"
	ViewScopeLineColor                   = "#ACADA4"
	ViewSelectScopeLineColor             = "#0BD01A"
	ViewScopePlayerColor                 = "#DA49D3"
	ViewSelectScopeLineWidthMin  float64 = 50
	ViewSelectScopeLineHeightMin float64 = 50
	ViewSelectScopeLineWidth     float64 = 200
	ViewSelectScopeLineHeight    float64 = 200
	ViewSelectScopeLineWidthMax  float64 = 2000
	ViewSelectScopeLineHeightMax float64 = 2000
	ViewMouseWheelSensitivity            = 3
	ViewFontColor1                       = "#0BD01A"
	ViewFontColor2                       = "#0BD01A"
	ViewFontSize                 float64 = 15
	ViewFontLineHeight           float64 = 20
)
