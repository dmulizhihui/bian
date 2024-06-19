// 主要做的工作是：
// 1，每次程序运行，都更新example.txt的单币信息；
// 2，当天第一次运行程序，更新A-history，记录当天的单币信息一次；
package persistentdata

import (
	"bufio"
	"fmt"
	"getSomething/environment"
	"os"
	"strconv"
	"strings"
)

func Persistentdata(accounts []environment.Coins) {
	text := ""
	// 打开原始文件以读写模式
	originalFile, err := os.OpenFile("example.txt", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer originalFile.Close()

	// 创建一个新文件以写入模式
	tempFile, err := os.CreateTemp("persistentdata", "example_temp.txt")
	if err != nil {
		fmt.Println("创建临时文件失败:", err)
		return
	}
	defer tempFile.Close()

	// 使用 bufio.Scanner 逐行读取原始文件内容
	scanner := bufio.NewScanner(originalFile)
	title := fmt.Sprintf("%-10s%-11s%-11s%-10s%-10s%-10s", "name", "oriPrice", "PerPrice", "CurPrice", "priceDiff", "diffPercent")
	fmt.Fprintln(tempFile, title)
	fmt.Println(title)
	cutLine := 0
	for scanner.Scan() {
		if cutLine == 0 {
			cutLine++
			continue
		}
		line := scanner.Text()
		lineTemp := strings.Fields(line)
		lineTemp = lineTemp[0:2]
		text += line
		text += "\n"
		// 写入文件每个coin的总价
		for _, account := range accounts {
			// 对于需要修改的行，进行数据追加
			if lineTemp[0] == account.Name && len(lineTemp) < 6 {
				sumPercent := 0.0
				priceDiff := 0.0
				originPrice, _ := strconv.ParseFloat(lineTemp[1], 64)
				priceDiff = account.SumPrices - originPrice
				sumPercent = (priceDiff / originPrice) * 100
				lineTemp = append(lineTemp, []string{
					fmt.Sprintf("%.10f ", account.PerPrices),
					fmt.Sprintf("%.2f", account.SumPrices),
					fmt.Sprintf("%.2f", priceDiff),
					fmt.Sprintf("|  %.2f%%  ", sumPercent),
				}...)

			}

		}
		temp := ""
		cut := 0
		for i := 0; i < len(lineTemp); i++ {
			switch cut {
			case 0:
				temp += fmt.Sprintf("%-10s", lineTemp[i])
			case 1:
				temp += fmt.Sprintf("%-10.8s", lineTemp[i])
			case 2:
				temp += fmt.Sprintf("%-11.10s", lineTemp[i])
			case 3:
				temp += fmt.Sprintf("%-10.8s", lineTemp[i])
			case 4:
				temp += fmt.Sprintf("%-10s", lineTemp[i])
			case 5:
				temp += fmt.Sprintf("%-6s", lineTemp[i])
			}
			cut++
		}
		//// 获取最大字符串的宽度
		// 将处理后的行写入临时文件
		_, err := fmt.Fprintln(tempFile, temp)
		fmt.Println(temp)
		if err != nil {
			fmt.Println("写入临时文件失败:", err)
			return
		}
		cutLine++
	}

	// 检查扫描过程中是否有错误
	if err := scanner.Err(); err != nil {
		fmt.Println("扫描文件失败:", err)
		return
	}

	// 关闭原始文件
	originalFile.Close()

	// 关闭临时文件
	tempFile.Close()

	// 移除原始文件
	err = os.Remove("example.txt")
	if err != nil {
		fmt.Println("移除原始文件失败:", err)
		return
	}

	// 重命名临时文件为原始文件
	err = os.Rename(tempFile.Name(), "example.txt")
	if err != nil {
		fmt.Println("重命名临时文件失败:", err)
		return
	}

	// fmt.Println("单币信息-已更新到 exmaple.txt")

}
