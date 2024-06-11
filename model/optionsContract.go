package model

import "time"

type OptionType string
type PositionType string

// Constants representing the possible values for OptionType and PositionType.
const (
	Call  OptionType   = "call"
	Put   OptionType   = "put"
	Long  PositionType = "long"
	Short PositionType = "short"
)

type OptionsContract struct {
	Type           OptionType   `json:"type"`
	StrikePrice    float64      `json:"strike_price"`
	Bid            float64      `json:"bid"`
	Ask            float64      `json:"ask"`
	ExpirationDate time.Time    `json:"expiration_date"`
	LongShort      PositionType `json:"long_short"`
}
