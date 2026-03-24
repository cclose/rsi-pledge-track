package model

type PledgeData struct {
	ID        int    `json:"ID" db:"id"`
	TimeStamp string `json:"TimeStamp" db:"pledge_timestamp"`
	Funds     int64  `json:"Funding" db:"funding"`
	Citizens  int32  `json:"Citizens" db:"citizens"`
	Fleet     int32  `json:"Fleet" db:"fleet"`
}
