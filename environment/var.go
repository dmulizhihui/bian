package environment

const (
	Url = "https://api.binance.com"
	//Url         = "https://199.59.149.236:443"
	BaApiKey    = "9gMUWEmkTqWGLCWsrFWEOooUF8MQu2ntCNIiqiazI4q9WCzLOrFn34dQ6PpE2FJM"
	BaApiSecret = "5BeOs2kLIEF5p7EegM4uEedlUGKrYD7I6kVvS02UmbKh8TiN4qvGlLYNqVZgR3bf"
	Signature   = "7316409891bbfe7599164e8d749cba5419b00f40f83432fba1c2cd2e29f0d1a5"
)

// 自定义切片排序
type Coins struct {
	Name       string
	Numbers    float64
	PerPrices  float64
	SumPrices  float64
	UpdateTime uint64
	//kLines    []*binance_connector.KlinesResponse
}

type SortCoins []Coins

func (s SortCoins) Len() int           { return len(s) }
func (s SortCoins) Less(i, j int) bool { return s[i].SumPrices > s[j].SumPrices }
func (s SortCoins) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
