package rekordbox

import (
	"fmt"
	"strings"
)

func GetHotCueColorName(pm PositionMark) string {
	rgb := fmt.Sprintf("%02X%02X%02X", pm.Red, pm.Green, pm.Blue)
	switch rgb {
	case "E62828":
		return "red"
	case "DE44CF":
		return "hotpink"
	case "FFFF00", "B4BE04", "C3AF04":
		return "yellow"
	case "28E214", "10B176":
		return "green"
	case "00E0FF", "50B4FF":
		return "aqua"
	case "305AFF", "6473FF":
		return "blue"
	case "B432FF", "AA72FF":
		return "purple"
	case "E0641B", "FFA500":
		return "orange"
	}
	return ""
}

func GetTrackColorName(hex string) string {
	switch strings.ToUpper(hex) {
	case "0XFF007F":
		return "pink"
	case "0XFF0000":
		return "red"
	case "0XFFA500":
		return "orange"
	case "0XFFFF00":
		return "yellow"
	case "0X00FF00":
		return "green"
	case "0X25FDE9":
		return "aqua"
	case "0X0000FF":
		return "blue"
	case "0X660099":
		return "purple"
	}
	return hex
}
