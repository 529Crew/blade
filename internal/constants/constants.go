package constants

import (
	"github.com/529Crew/blade/internal/config"
	"github.com/gagliardetto/solana-go"
)

var (
	Wallet         = solana.MustPrivateKeyFromBase58(config.Get().WalletPrivateKey)
	DummySignature = solana.MustSignatureFromBase58("1111111111111111111111111111111111111111111111111111111111111111")
)
