package types

import "github.com/gagliardetto/solana-go"

type MarketStateV2 struct {
	Reserved               [5]byte `json:"-"`
	AccountFlags           uint64
	OwnAddress             solana.PublicKey
	VaultSignerNonce       uint64
	BaseMint               solana.PublicKey
	QuoteMint              solana.PublicKey
	BaseVault              solana.PublicKey
	BaseDepositsTotal      uint64
	BaseFeesAccrued        uint64
	QuoteVault             solana.PublicKey
	QuoteDepositsTotal     uint64
	QuoteFeesAccrued       uint64
	QuoteDustThreshold     uint64
	RequestQueue           solana.PublicKey
	EventQueue             solana.PublicKey
	Bids                   solana.PublicKey
	Asks                   solana.PublicKey
	BaseLotSize            uint64
	QuoteLotSize           uint64
	FeeRateBps             uint64
	ReferrerRebatesAccrued uint64
	Padding                [7]byte `json:"-"`
}

type MarketStateV3 struct {
	Reserved               [5]uint8
	AccountFlags           uint64
	OwnAddress             solana.PublicKey
	VaultSignerNonce       uint64
	BaseMint               solana.PublicKey
	QuoteMint              solana.PublicKey
	BaseVault              solana.PublicKey
	BaseDepositsTotal      uint64
	BaseFeesAccrued        uint64
	QuoteVault             solana.PublicKey
	QuoteDepositsTotal     uint64
	QuoteFeesAccrued       uint64
	QuoteDustThreshold     uint64
	RequestQueue           solana.PublicKey
	EventQueue             solana.PublicKey
	Bids                   solana.PublicKey
	Asks                   solana.PublicKey
	BaseLotSize            uint64
	QuoteLotSize           uint64
	FeeRateBps             uint64
	ReferrerRebatesAccrued uint64
	Authority              solana.PublicKey
	PruneAuthority         solana.PublicKey
	ConsumeEventsAuthority solana.PublicKey
	Reserved2              [992]uint8
	Padding                [7]uint8
}
