package paint

import (
	"bufio"
	"fmt"
	consts "getSomething/const"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"

	// "gonum.org/v1/plot/plotutil"
	"image/color" // 确保导入了这个包
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot/vg"
)

func Mkdir(text string) {
	// 工程的A目录路径，请根据实际情况进行调整
	dirPath := consts.HistoryCoinDataDir
	// 目标文件名

	// ### 加载北京时区 ###
	loc, err := time.LoadLocation("Asia/Shanghai") // Beijing is in the Asia/Shanghai time zone
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	// 获取当前时间并转换为北京时区
	now := time.Now().In(loc)
	// 转换时间到字符串
	// 你可以根据需要调整时间格式
	// 这里使用的格式是：年-月-日 时:分:秒
	curTime := now.Format("2006-01-02")
	filename := "z-" + curTime + ".txt"

	// 确保目录存在
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Printf("failed to create directory: %s", err)
	}
	// 在A目录下尝试创建文件，仅当它不存在时
	filePath := filepath.Join(dirPath, filename)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsExist(err) {
			// 文件已存在的情况
			fmt.Printf("file %s already exists \n", filePath)
			return
		}
		// 处理其他可能的错误
		fmt.Printf("error creating file: %s", err)
	}
	//   ##########  文件第一次创建的时候 会写入信息    ###########
	if _, err := file.Write([]byte(text)); err != nil {
		fmt.Println("写入文件时发生错误:", err)
		return
	}
	// fmt.Println("成功写入文件。")
	defer file.Close()
	// 文件创建成功
	// fmt.Printf("file %s created successfully in %s directory", filename, dirPath)
}

func Paint(dataName string) {
	fileName := consts.HistoryCoinDataDir + "a-" + dataName + ".txt" // 数据文件的名称
	pts, err := readData(fileName)
	if err != nil {
		fmt.Printf("读取或解析数据时出错: %v\n", err)
		return
	}
	// 创建一个plot实例
	p := plot.New()

	p.Title.Text = "折线图"
	p.X.Label.Text = "日期"
	p.Y.Label.Text = "数值"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02 15:04:05"}
	// 添加折线图
	// 创建一条线，并设置为灰色
	line, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	line.Color = color.Gray{Y: 100} // Y: 0~255, 定义灰度等级

	// 创建数据点 为灰色
	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Color = color.Gray{Y: 100} // 设置数据点为灰色
	scatter.GlyphStyle.Radius = vg.Points(1)      // 可以调整数据点大小
	// 添加线到plot
	p.Add(line, scatter)

	// 保存为PNG图像

	if err := p.Save(10*vg.Inch, 10*vg.Inch, consts.HistoryCoinDataPlotDir+"a-"+dataName+"-plot.png"); err != nil {
		panic(err)
	}
	// 关闭  折线生成完毕的 日志
	// fmt.Println(dataName + " 折线图已生成")
}

// readData 从给定的文件名读取数据，并返回一个plotter.XYs类型的数据点切片
func readData(fileName string) (plotter.XYs, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// 忽略最后一个空白行（如果有）
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	var pts plotter.XYs
	for _, line := range lines {
		parts := strings.Split(line, "&&")
		if len(parts) != 2 {
			continue // 跳过格式不正确的行
		}
		valueStr := strings.TrimSpace(parts[0])
		dateStr := strings.TrimSpace(parts[1])
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("解析数值'%s'失败: %w", valueStr, err)
		}
		date, err := time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return nil, fmt.Errorf("解析日期'%s'失败: %w", dateStr, err)
		}
		// 由于gonum/plot无法直接处理时间类型，我们将时间转换为float64
		// 这里简单地使用时间的Unix时间戳
		pts = append(pts, plotter.XY{X: float64(date.Unix()), Y: value})
	}
	return pts, nil
}
