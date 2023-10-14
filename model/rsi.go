package model

type ChartType string

const ChartTypeDay ChartType = "day"

type RSICrowdFundStatsRequest struct {
	Chart ChartType `json:"chart"`
	Fleet bool      `json:"fleet"`
	Fans  bool      `json:"fans"`
	Funds bool      `json:"funds"`
}

type RSIResponse struct {
	Success int         `json:"success"`
	Code    string      `json:"code"`
	Msg     string      `json:"msg"`
	Data    DataSection `json:"data"`
}

type DataSection struct {
	Chart ChartData `json:"chart"`
	Fans  int       `json:"fans"`
	Funds int64     `json:"funds"`
	Fleet string    `json:"fleet"`
}

// The chart data can be represented as a map with dates as keys
type ChartData map[string]ChartEntry

type ChartEntry struct {
	Gross string `json:"gross"`
	Axis  string `json:"axis"`
}
