package mathtool

import "math"

const epsilon = 1e-9 // 定义一个小的误差范围

// 比较两个浮点数是否相等
func FloatEquals(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

// 比较两个浮点数的大小关系
func CompareFloats(a, b float64) bool {
	if FloatEquals(a, b) {
		return true // 相等
	} else if a < b {
		return true // 小于
	} else {
		return false // 大于
	}
}
