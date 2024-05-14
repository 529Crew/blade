package pump_tx

import (
	"context"
	"math/rand"

	"github.com/529Crew/blade/idls/pump"
	"github.com/529Crew/blade/internal/client"
	"github.com/529Crew/blade/internal/config"
	"github.com/529Crew/blade/internal/constants"
	"github.com/529Crew/blade/internal/logger"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
)

func Buy(buyInst *pump.Buy) error {
	instructions := []solana.Instruction{
		computebudget.NewSetComputeUnitLimitInstruction(200_420).Build(),
		computebudget.NewSetComputeUnitPriceInstruction(config.Get().BuyPriorityFee + uint64(rand.Intn(1000))).Build(),
	}

	mint := buyInst.GetMintAccount().PublicKey.String()

	logger.Log.Printf("attempting to buy %s\n", mint)

	associatedUser, _, err := solana.FindAssociatedTokenAddress(constants.Wallet.PublicKey(), buyInst.GetMintAccount().PublicKey)
	if err != nil {
		logger.Log.Printf("error buying %s: %v", mint, err)
		return err
	}

	_, err = client.GetUtil().GetAccountInfo(context.Background(), associatedUser)
	if err != nil {
		if err.Error() == "not found" {
			instructions = append(instructions,
				associatedtokenaccount.NewCreateInstruction(
					constants.Wallet.PublicKey(),
					constants.Wallet.PublicKey(),
					buyInst.GetMintAccount().PublicKey,
				).Build(),
			)
		} else {
			logger.Log.Printf("error buying %s: %v", mint, err)
			return err
		}
	}

	instructions = append(instructions,
		pump.NewBuyInstruction(
			5_000_000,        /* need to figure out a way to properly calculate quote here */
			1_000_000_000/10, /* 0.1 SOL */
			buyInst.GetGlobalAccount().PublicKey,
			buyInst.GetFeeRecipientAccount().PublicKey,
			buyInst.GetMintAccount().PublicKey,
			buyInst.GetBondingCurveAccount().PublicKey,
			buyInst.GetAssociatedBondingCurveAccount().PublicKey,
			associatedUser,
			constants.Wallet.PublicKey(),
			buyInst.GetSystemProgramAccount().PublicKey,
			buyInst.GetTokenProgramAccount().PublicKey,
			buyInst.GetRentAccount().PublicKey,
			buyInst.GetEventAuthorityAccount().PublicKey,
			buyInst.GetProgramAccount().PublicKey,
		).Build(),
	)

	block, err := client.Get().GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		logger.Log.Printf("error buying %s: %v", mint, err)
		return err
	}

	tx, err := solana.NewTransaction(
		instructions,
		block.Value.Blockhash,
		solana.TransactionPayer(constants.Wallet.PublicKey()),
	)
	if err != nil {
		logger.Log.Printf("error buying %s: %v", mint, err)
		return err
	}

	if _, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if constants.Wallet.PublicKey().Equals(key) {
				return &constants.Wallet
			}
			return nil
		},
	); err != nil {
		logger.Log.Printf("error buying %s: %v", mint, err)
		return err
	}

	logger.Log.Printf("submitting tx to buy %s\n", mint)
	sig, err := client.Get().SendTransactionWithOpts(context.Background(), tx, rpc.TransactionOpts{
		SkipPreflight: true,
	})
	logger.Log.Printf("submitted tx to buy %s: %s\n", mint, sig)

	return nil
}
