package tools

import (
	"fmt"
	"time"
)

func ToBj(nowTime uint64) {
	// 以毫秒为单位的时间戳
	milliseconds := int64(nowTime)

	// 将毫秒转换为秒
	seconds := milliseconds / 1000

	// 将秒转换为时间对象
	t := time.Unix(seconds, 0).UTC().Add(8 * time.Hour) // 将时间调整为北京时间

	// 格式化输出时间
	fmt.Println("Beijing Time:", t.Format("2006-01-02 15:04:05"))
}
