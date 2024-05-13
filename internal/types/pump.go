package types

import "github.com/gagliardetto/solana-go"

type BondingCurve struct {
	VirtualTokenReserves uint64
	VirtualSolReserves   uint64
	RealTokenReserves    uint64
	RealSolReserves      uint64
	TokenTotalSupply     uint64
	Complete             bool
}

type Global struct {
	Initialized                 bool
	Authority                   solana.PublicKey
	FeeRecipient                solana.PublicKey
	InitialVirtualTokenReserves uint64
	InitialVirtualSolReserves   uint64
	InitialRealTokenReserves    uint64
	TokenTotalSupply            uint64
	FeeBasisPoints              uint64
}

type Coins []struct {
	Mint                   string  `json:"mint"`
	Name                   string  `json:"name"`
	Symbol                 string  `json:"symbol"`
	Description            string  `json:"description"`
	ImageURI               string  `json:"image_uri"`
	MetadataURI            string  `json:"metadata_uri"`
	Twitter                string  `json:"twitter"`
	Telegram               string  `json:"telegram"`
	BondingCurve           string  `json:"bonding_curve"`
	AssociatedBondingCurve string  `json:"associated_bonding_curve"`
	Creator                string  `json:"creator"`
	CreatedTimestamp       int64   `json:"created_timestamp"`
	RaydiumPool            string  `json:"raydium_pool"`
	Complete               bool    `json:"complete"`
	VirtualSolReserves     int64   `json:"virtual_sol_reserves"`
	VirtualTokenReserves   int64   `json:"virtual_token_reserves"`
	Hidden                 any     `json:"hidden"`
	TotalSupply            int64   `json:"total_supply"`
	Website                string  `json:"website"`
	ShowName               bool    `json:"show_name"`
	LastTradeTimestamp     int64   `json:"last_trade_timestamp"`
	KingOfTheHillTimestamp int64   `json:"king_of_the_hill_timestamp"`
	MarketCap              float64 `json:"market_cap"`
	ReplyCount             int     `json:"reply_count"`
	LastReply              int64   `json:"last_reply"`
	Nsfw                   bool    `json:"nsfw"`
	MarketID               string  `json:"market_id"`
	Inverted               bool    `json:"inverted"`
	Username               string  `json:"username"`
	ProfileImage           any     `json:"profile_image"`
	UsdMarketCap           float64 `json:"usd_market_cap"`
}
