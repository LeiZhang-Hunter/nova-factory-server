package math

import (
	"math"
)

// RoundFloat 四舍五入到指定的小数位数
func RoundFloat(num float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(num*ratio) / ratio
}
