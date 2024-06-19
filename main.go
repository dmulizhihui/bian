package main

import (
	"context"
	"fmt"
	consts "getSomething/const"
	"getSomething/mathtool"
	"getSomething/paint"
	"getSomething/persistentdata"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	//ba_client "getSomething/client"
	env "getSomething/environment"
	binance_connector "github.com/binance/binance-connector-go"
)

var (
	apiKey    string = env.BaApiKey
	secretKey string = env.BaApiSecret
	baseURL   string = env.Url
)

func getAllKindOfCoinsAndNullCoins(client *binance_connector.Client) ([]string, []string) {
	allCoinsInfo, err := client.NewGetAllCoinsInfoService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return []string{}, []string{}
	}
	coinsName := []string{}
	assetCoins := []string{}
	for _, coinsInfo := range allCoinsInfo {
		coinsName = append(coinsName, coinsInfo.Coin)
		if coinsInfo.Free != "0" {
			assetCoins = append(assetCoins, coinsInfo.Coin)
			fmt.Println(coinsInfo.Coin, coinsInfo.Free)
		}
	}
	return coinsName, assetCoins
}
func getAccountsInfo(client *binance_connector.Client, coin string) string {
	// 获取现货账户余额
	res := ""
	accountResponse, _ := client.NewGetAccountService().Do(context.Background())
	for _, balance := range accountResponse.Balances {
		if coin != "" {
			if coin == balance.Asset {
				res = balance.Free
				fmt.Printf("现货_账户余额 %s: %s\n", balance.Asset, balance.Free)
			}
		} else {
			// 获取现货账户余额，USDT
			//if balance.Asset == consts.USDT_CoinName {
			//	fmt.Printf("现货_账户余额 %s: %s\n", consts.USDT_CoinName, balance.Free)
			//}
			// 获取现货资金不为0 的所有资产
			f64, _ := strconv.ParseFloat(balance.Free, 64)
			if mathtool.FloatEquals(f64, 0.0) == false {
				fmt.Printf("现货_账户余额 %s: %s\n", balance.Asset, balance.Free)
			}
		}
	}
	return res
}
func strategyOfBuy(client *binance_connector.Client, coinsSlic []string, remainUSDT string) {
	//usdt, _ := strconv.ParseInt(remainUSDT, 10, 64)
	for _, coinName := range coinsSlic {
		// 获取coins的 K 线，判断最近3min 都是红的
		sysbol := string("REI") + "USDT"

		klines := getKLinesOfCoin(client, sysbol)

		//if err != nil {
		//	fmt.Println(err)
		//}
		// 标记要每次都要上涨
		up := true
		for _, kline := range klines {
			openPrice, _ := strconv.ParseFloat(kline.Open, 64)
			closePrice, _ := strconv.ParseFloat(kline.Close, 64)
			//判断最近3min是不是都在上涨
			if !mathtool.CompareFloats(openPrice, closePrice) {
				up = false
				break
			}
		}
		if up {
			fmt.Println(coinName)
		}
		//fmt.Println(binance_connector.PrettyPrint(klines))
	}
}

func getKLinesOfCoin(client *binance_connector.Client, sysbol string) []*binance_connector.KlinesResponse {
	klines, err := client.NewKlinesService().Symbol(sysbol + "USDT").Interval("30s").Limit(10).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return klines
}

func main() {
	// Initialise the client
	client := binance_connector.NewClient(apiKey, secretKey, baseURL)

	// 1 获取账户中 现存的币名字 及数量 ,及价值
	accounts := []env.Coins{}

	// Binance Account Information (USER_DATA) - GET /api/v3/account
	accountInformation, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

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
	curTime := now.Format("2006-01-02 15:04:05")

	err = os.MkdirAll(consts.HistoryCoinDataDir, 0755)
	if err != nil {
		fmt.Printf("failed to create directory: %s", err)
	}

	for _, balance := range accountInformation.Balances {
		if balance.Free != "0.00000000" && balance.Free != "0.00" && balance.Free != "0.0" {
			coinName := strings.TrimPrefix(balance.Asset, "LD")
			numberOfCoin, _ := strconv.ParseFloat(balance.Free, 32)
			if coinName == "USDT" {
				continue
			}
			// 获取实时价格
			tickerPrice, err := client.NewTickerPriceService().Symbol(coinName + "USDT").Do(context.Background())
			if err != nil {
				fmt.Println(coinName, err)
				return
			}
			currentPrice, _ := strconv.ParseFloat(tickerPrice.Price, 32)
			accounts = append(accounts, env.Coins{
				Name:      coinName,
				Numbers:   numberOfCoin,
				PerPrices: currentPrice,
				// 计算出总价
				SumPrices:  numberOfCoin * currentPrice * consts.ExchangeRate,
				UpdateTime: accountInformation.UpdateTime,
			})

			// ##### 落盘历史总价 到文件 #######
			fileName := consts.HistoryCoinDataDir + "a-" + coinName + ".txt"
			// 以只写方式打开文件，如果文件不存在，则创建该文件
			// os.O_CREATE | os.O_WRONLY 表示如果不存在则创建，且为只写模式
			// 0666 是文件权限设置，表示文件所有者、所属组和其他人都可以读写该文件
			file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				fmt.Println("打开文件出错:", err)
				return
			}
			defer file.Close() // 确保在函数退出时关闭文件
			// 定义要写入的数据
			// 将数据写入文件
			sumStr := fmt.Sprintf("%.6f", numberOfCoin*currentPrice*consts.ExchangeRate)
			if _, err := file.WriteString(sumStr + "  &&  " + curTime + "\n"); err != nil {
				return
			}
			// fmt.Printf("%s 数据写入成功 \n", coinName)
			paint.Paint(coinName)
		}
	}

	sum := 0.0
	cut := 0
	sort.Sort(env.SortCoins(accounts))
	fmt.Println("\n当前北京时间：", curTime)
	fmt.Println()
	for _, account := range accounts {
		if mathtool.CompareFloats(account.SumPrices, 1000) {
			continue
		}
		fmt.Println(fmt.Sprintf("%-13s%-13f%-13f", account.Name, account.SumPrices, account.PerPrices))
	}
	for _, account := range accounts {
		sum += account.SumPrices
		// fmt.Println(account.Name, account.SumPrices, "===========", sum)
		// 只统计 总价大于 5RMB
		if mathtool.CompareFloats(account.SumPrices, 5000.0) {
			continue
		}
		cut++
	}
	fmt.Println("上述种类:", cut)
	fmt.Println()

	// 在这里 做一手，，挂的卖单 不在sum里面显示 ，额外加上这个差价
	// sum += consts.AddGapOfSaleing

	persistentdata.Persistentdata(accounts)

	btc, _ := client.NewTickerPriceService().Symbol("BTCUSDT").Do(context.Background())

	fmt.Println("\nbtc:", btc.Price)
	fmt.Println("sum:", sum)

	// ##### 落盘历史总价 到文件 #######
	fileName := consts.HistoryCoinDataDir + "a-sum.txt"
	// 以只写方式打开文件，如果文件不存在，则创建该文件
	// os.O_CREATE | os.O_WRONLY 表示如果不存在则创建，且为只写模式
	// 0666 是文件权限设置，表示文件所有者、所属组和其他人都可以读写该文件
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("打开文件出错:", err)
		return
	}
	defer file.Close() // 确保在函数退出时关闭文件
	// 定义要写入的数据
	data := sum
	// 将数据写入文件
	sumStr := fmt.Sprintf("%.6f", data)
	if _, err := file.WriteString(sumStr + "  &&  " + curTime + "\n"); err != nil {
		return
	}
	// fmt.Println("a-sum 数据写入成功")
	paint.Paint("sum")

	// ##### 落盘历史总价 到文件 #######
	fileNameBTC := consts.HistoryCoinDataDir + "a-btc.txt"
	// 以只写方式打开文件，如果文件不存在，则创建该文件
	// os.O_CREATE | os.O_WRONLY 表示如果不存在则创建，且为只写模式
	// 0666 是文件权限设置，表示文件所有者、所属组和其他人都可以读写该文件
	fileBTC, err := os.OpenFile(fileNameBTC, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("打开文件出错:", err)
		return
	}
	defer file.Close() // 确保在函数退出时关闭文件
	// 定义要写入的数据
	data1 := btc.Price + "  &&  " + curTime
	// 将数据写入文件
	if _, err := fileBTC.WriteString(data1 + "\n"); err != nil {
		return
	}
	// fmt.Println("a-btc数据写入成功")
	paint.Paint("btc")

	/////////把 btc 和 sum 单独写到一个文件 ////////////
	//  把btc 和 sum 单独放到一个文件，并且是覆盖更新；
	// 打开原始文件以读写模式
	//   ***  如下这句话 不要乱用，会清空文件的 *********
	btcSumFile, err := os.OpenFile("z--btc-sum.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer btcSumFile.Close()
	btcSumdate := fmt.Sprintln("北京时间:", curTime)
	btcSumdate += fmt.Sprintln("  ")
	btcSumdate += fmt.Sprintln("btc:", btc.Price)
	btcSumdate += fmt.Sprintln("  ")
	btcSumdate += fmt.Sprintln("sum:", sumStr)

	_, err = fmt.Fprintln(btcSumFile, btcSumdate)
	if err != nil {
		fmt.Println("btcSumdate写入文件失败:", err)
		return
	}
	/////////////////////

}
