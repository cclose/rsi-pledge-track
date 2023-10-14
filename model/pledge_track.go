package model

type PledgeData struct {
	ID        int    `json:"ID" db:"id"`
	TimeStamp string `json:"TimeStamp" db:"timestamp"`
	Funds     int64  `json:"Funding" db:"funding"`
	Citizens  int    `json:"Citizens" db:"citizens"`
	Fleet     int    `json:"Fleet" db:"fleet"`
}
