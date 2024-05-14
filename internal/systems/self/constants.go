package self_system

import (
	"github.com/529Crew/blade/internal/config"
	"github.com/gagliardetto/solana-go"
)

var (
	SELF = solana.MustPrivateKeyFromBase58(config.Get().WalletPrivateKey).PublicKey()
)
